package generator

import (
	"context"
	"fmt"
	"load-generation-system/internal/core"
	"load-generation-system/internal/scenarios"
	"load-generation-system/internal/service/callers"
	"load-generation-system/internal/service/http"
	"load-generation-system/pkg/scheduler"
	"log"
	"sync"
	"time"
)

// increment represents a group of users executing a specific operation within an attack
type increment struct {
	operationID string             // Unique identifier for the operation
	users       []*user            // Collection of virtual users in this increment
	ctx         context.Context    // Context for managing increment lifecycle
	cancel      context.CancelFunc // Function to cancel the increment
}

// attack represents a complete load test consisting of multiple increments
type attack struct {
	increments map[int64]increment // Map of increments by their IDs
	ctx        context.Context     // Context for managing attack lifecycle
	cancel     context.CancelFunc  // Function to cancel the attack
	jobID      string              // Scheduler job identifier
}

// Config contains configuration parameters for the load generator
type Config struct {
	UsersPerClient        int64 // Number of users sharing a single HTTP client
	MinIdleConnTimeoutSec int64 // Minimum idle connection timeout in seconds
	MaxIdleConnTimeoutSec int64 // Maximum idle connection timeout in seconds
}

// generator is the main implementation of the LoadGenerator interface
type generator struct {
	attacks   map[int64]attack     // Active attacks indexed by attack ID
	mu        sync.RWMutex         // Mutex for concurrent access to attacks
	ctx       context.Context      // Root context for the generator
	cancel    context.CancelFunc   // Function to shutdown the generator
	stop      sync.WaitGroup       // WaitGroup for graceful shutdown
	scheduler *scheduler.Scheduler // Job scheduler for attack execution
	config    Config               // Generator configuration
}

func New(config Config) core.LoadGenerator {
	ctx, cancel := context.WithCancel(context.Background())

	return &generator{
		attacks:   make(map[int64]attack),
		ctx:       ctx,
		cancel:    cancel,
		scheduler: scheduler.New(),
		config:    config,
	}
}

// StartAttack initiates a new attack or adds an increment to an existing attack
//
// Parameters:
//   - start: Operation details including attack ID, increment ID, and scenarios
//
// Returns:
//   - error: Any error that occurs during attack initialization
func (g *generator) StartAttack(start core.OperationStart) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	att, exists := g.attacks[start.AttackID]
	if !exists {
		// Schedule new attack if it doesn't exist
		jobID, err := g.scheduler.NewJob(
			start.WaitTimeSec,
			func() {
				g.executeAttack(start.AttackID, start.WaitTimeSec)
			},
		)
		if err != nil {
			return fmt.Errorf("failed to schedule attack %d: %v", start.AttackID, err)
		}

		ctx, cancel := context.WithCancel(g.ctx)
		att = attack{
			increments: make(map[int64]increment),
			ctx:        ctx,
			cancel:     cancel,
			jobID:      jobID,
		}
		g.attacks[start.AttackID] = att
	}

	// Check for duplicate increment
	if increment, exists := att.increments[start.IncrementID]; exists && increment.operationID == start.ID {
		return fmt.Errorf(
			"attack %d, increment %d within the operation %s has already been created",
			start.AttackID, start.IncrementID, start.ID,
		)
	}

	// Create users for each scenario
	var users []*user
	var httpClient core.Client

	for name, count := range start.Scenarios {
		scenario, ok := scenarios.AvailableScenarios[name]
		if !ok {
			log.Printf("scenario %s is not existed! It will be skipped", name)
			continue
		}

		for i := int64(0); i < count; i++ {
			// Create new HTTP client when needed
			if i%g.config.UsersPerClient == 0 {
				httpClient = http.NewClient(
					g.config.MinIdleConnTimeoutSec,
					g.config.MaxIdleConnTimeoutSec,
				)
			}

			caller := callers.NewCaller(httpClient)
			users = append(users, newUser(fmt.Sprintf("user for %s #%d", name, i), scenario, caller))
			g.stop.Add(1)
		}
	}

	// Create and store the new increment
	ctx, cancel := context.WithCancel(att.ctx)
	att.increments[start.IncrementID] = increment{
		operationID: start.ID,
		users:       users,
		ctx:         ctx,
		cancel:      cancel,
	}

	return nil
}

// executeAttack coordinates the execution of all users in an attack with proper pacing
//
// Parameters:
//   - attackID: Identifier of the attack to execute
//   - waitTimeSec: Total time over which to distribute user starts
func (g *generator) executeAttack(attackID int64, waitTimeSec float64) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	att, exists := g.attacks[attackID]
	if !exists {
		return
	}

	// Calculate total users and start interval
	var userCounter int
	for _, increment := range att.increments {
		userCounter += len(increment.users)
	}
	if userCounter == 0 {
		return
	}

	interval := time.Duration(float64(waitTimeSec) / float64(userCounter) * float64(time.Second))

	// Start users with calculated interval
	for _, increment := range att.increments {
		for _, user := range increment.users {
			go user.Run(increment.ctx)
			time.Sleep(interval)
		}
	}
}

// StopAttack terminates either a specific increment or an entire attack
//
// Parameters:
//   - stop: Operation details including what to stop
//
// Returns:
//   - error: Any error that occurs during attack termination
func (g *generator) StopAttack(stop core.OperationStop) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	attack, exists := g.attacks[stop.AttackID]
	if !exists {
		return core.ErrAttackNotFound
	}

	if stop.IncrementID != nil {
		// Stop specific increment
		increments := g.attacks[stop.AttackID].increments
		increment, exists := increments[*stop.IncrementID]
		if !exists {
			return core.ErrAttackNotFound
		}

		increment.cancel()
		delete(increments, *stop.IncrementID)

		// Clean up users
		for _, user := range increment.users {
			go func() {
				defer g.stop.Done()
				user.Destroy(g.ctx)
			}()
		}

		// Clean up attack if no increments remain
		if len(increments) == 0 {
			attack.cancel()
			if err := g.scheduler.RemoveJob(g.attacks[stop.AttackID].jobID); err != nil {
				return fmt.Errorf("unable to remove attack job: %v", err)
			}
			delete(g.attacks, stop.AttackID)
		}
	} else {
		// Stop entire attack
		attack.cancel()
		if err := g.scheduler.RemoveJob(g.attacks[stop.AttackID].jobID); err != nil {
			return fmt.Errorf("unable to remove attack job: %v", err)
		}
		delete(g.attacks, stop.AttackID)

		// Clean up all users
		for _, increment := range attack.increments {
			for _, user := range increment.users {
				go func() {
					defer g.stop.Done()
					user.Destroy(g.ctx)
				}()
			}
		}
	}

	return nil
}

// Stop gracefully shuts down the generator, terminating all active attacks
func (g *generator) Stop() {
	ctx := context.Background()

	g.mu.Lock()
	defer g.mu.Unlock()

	g.cancel()
	if err := g.scheduler.Shutdown(); err != nil {
		log.Printf("error shutdowning scheduler: %v", err)
	}

	// Clean up all users in all attacks
	for _, attack := range g.attacks {
		for _, increment := range attack.increments {
			for _, user := range increment.users {
				go func() {
					defer g.stop.Done()
					user.Destroy(ctx)
				}()
			}
		}
	}

	g.stop.Wait()
}
