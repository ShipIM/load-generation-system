package node

import (
	"context"
	"load-generation-system/internal/core"
	"slices"
	"sync"
	"time"
)

// node implements the core.Node interface and represents a single node in the load generation system.
// It manages attack operations, maintains state, and handles operation retries.
type node struct {
	name      string                          // Name identifier for the node
	scenarios map[string]core.ScenarioDetails // Available scenarios on this node
	attacks   map[int64]core.AttackDetails    // Active attacks on this node
	isActive  bool                            // Whether the node is currently active

	ops     chan core.Operation // Channel for sending operations to the controller
	opQueue chan core.Operation // Buffered queue for pending operations
	ack     chan any            // Channel for operation acknowledgments

	retryInterval time.Duration // Interval between operation retry attempts

	mu sync.Mutex // Mutex to protect concurrent access to node state
}

func New(
	name string,
	scenarios map[string]core.ScenarioDetails,
	ops chan core.Operation,
	opQueueCapacity, retryIntervalSec int64,
) core.Node {
	return &node{
		name:          name,
		scenarios:     scenarios,
		attacks:       make(map[int64]core.AttackDetails),
		ops:           ops,
		opQueue:       make(chan core.Operation, opQueueCapacity),
		retryInterval: time.Duration(retryIntervalSec) * time.Second,
		ack:           make(chan any, 1),
	}
}

// Start activates the node and begins processing operations.
// It launches a goroutine to manage operation processing and retries.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
func (n *node) Start(ctx context.Context) {
	n.isActive = true
	go n.processOps(ctx) // Start operation processing in background
}

// processOps handles the core operation processing loop for the node.
// It manages operation retries and acknowledgments.
func (n *node) processOps(ctx context.Context) {
	var currentOp *core.Operation // Currently pending operation
	n.ack <- nil                  // Initialize acknowledgment channel

	timer := time.NewTimer(n.retryInterval)
	defer func() {
		timer.Stop()
		n.isActive = false // Mark as inactive when stopping
	}()

	for {
		select {
		case <-ctx.Done(): // Context cancellation
			return
		case <-timer.C: // Retry timeout
			n.mu.Lock()
			if currentOp != nil {
				n.ops <- *currentOp // Retry operation
				timer.Reset(n.retryInterval)
			}
			n.mu.Unlock()
		case <-n.ack: // Operation acknowledged
			n.mu.Lock()
			currentOp = nil // Clear current operation
			n.mu.Unlock()

			select {
			case <-ctx.Done():
				return
			case op := <-n.opQueue: // Get next operation
				n.mu.Lock()
				currentOp = &op
				n.ops <- *currentOp // Send new operation
				timer.Reset(n.retryInterval)
				n.mu.Unlock()
			}
		}
	}
}

// AckOperation signals that the current operation has been acknowledged.
func (n *node) AckOperation() {
	n.ack <- nil
}

// StartAttack initiates a new attack or updates an existing one.
// It validates scenarios and queues the start operation.
//
// Parameters:
//   - start: Operation details for starting the attack
//
// Returns:
//   - error: ErrScenarioNotFound if invalid scenario, nil otherwise
func (n *node) StartAttack(start core.OperationStart) error {
	n.mu.Lock()
	// Prepare increment details
	increment := core.IncrementDetails{
		ID:        start.IncrementID,
		AttackID:  start.AttackID,
		Scenarios: start.Scenarios,
	}

	// Update existing attack or create new one
	if attack, exists := n.attacks[start.AttackID]; exists {
		var incrementDetails *core.IncrementDetails
		// Find existing increment
		for _, increment := range attack.Increments {
			if increment.ID == start.IncrementID {
				temp := increment
				incrementDetails = &temp
				break
			}
		}

		if incrementDetails != nil {
			// Update existing increment counters
			for name, counter := range start.Scenarios {
				scenarioCounter := incrementDetails.Scenarios[name]
				scenarioCounter += counter
				incrementDetails.Scenarios[name] = scenarioCounter
			}
		} else {
			// Add new increment to existing attack
			attack.Increments = append(attack.Increments, increment)
			n.attacks[start.AttackID] = attack
		}
	} else {
		// Create new attack
		incrementDetails := []core.IncrementDetails{increment}
		n.attacks[start.AttackID] = core.AttackDetails{
			ID:         start.AttackID,
			Increments: incrementDetails,
		}
	}
	n.mu.Unlock()

	// Validate all scenarios exist
	for scenario := range start.Scenarios {
		if _, ok := n.scenarios[scenario]; !ok {
			n.opQueue <- core.Operation{
				Kill: &core.OperationKill{}, // Send kill if invalid scenario
			}
			return core.ErrScenarioNotFound
		}
	}

	// Queue start operation
	n.opQueue <- core.Operation{
		Start: &start,
	}

	return nil
}

// StopAttack terminates an existing attack or increment.
// It updates the attack state and queues the stop operation.
//
// Parameters:
//   - stop: Operation details for stopping the attack
//
// Returns:
//   - error: ErrAttackNotFound or ErrIncrementNotFound if invalid ID
func (n *node) StopAttack(stop core.OperationStop) error {
	n.mu.Lock()
	attack, exists := n.attacks[stop.AttackID]
	if !exists {
		n.mu.Unlock()
		return core.ErrAttackNotFound
	}

	// Handle increment-specific stop
	if stop.IncrementID != nil {
		var incrementID *int
		// Find the increment to remove
		for i, increment := range attack.Increments {
			if increment.ID == *stop.IncrementID {
				temp := i
				incrementID = &temp
				break
			}
		}
		if incrementID == nil {
			n.mu.Unlock()
			return core.ErrIncrementNotFound
		}

		// Remove increment or entire attack
		if len(attack.Increments) > 1 {
			attack.Increments = slices.Delete(attack.Increments, *incrementID, *incrementID+1)
			n.attacks[stop.AttackID] = attack
		} else {
			delete(n.attacks, stop.AttackID)
		}
	} else {
		// Remove entire attack
		delete(n.attacks, stop.AttackID)
	}
	n.mu.Unlock()

	// Queue stop operation
	n.opQueue <- core.Operation{
		Stop: &stop,
	}

	return nil
}

// GetDetails returns the current state and configuration of the node.
//
// Returns:
//   - core.NodeDetails: Snapshot of the node's current state
func (n *node) GetDetails() core.NodeDetails {
	n.mu.Lock()
	// Copy active attacks
	attacks := make([]core.AttackDetails, 0, len(n.attacks))
	for _, attack := range n.attacks {
		attacks = append(attacks, attack)
	}
	n.mu.Unlock()

	// Copy available scenarios
	scenarios := make([]core.ScenarioDetails, 0, len(n.scenarios))
	for _, scenario := range n.scenarios {
		scenarios = append(scenarios, scenario)
	}

	return core.NodeDetails{
		Name:      n.name,
		Scenarios: scenarios,
		Attacks:   attacks,
		IsActive:  n.isActive,
	}
}
