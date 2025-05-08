package generator

import (
	"context"
	"load-generation-system/internal/metrics"
	"load-generation-system/internal/scenarios"
	"load-generation-system/internal/service/callers"
	"log"
	"sync"
)

// user represents a virtual user that runs a scenario in a load generation system.
// Each user runs a scenario using a specified caller and is controlled using synchronization.
type user struct {
	name     string             // The name of the user.
	scenario scenarios.Scenario // The scenario this user is running.
	caller   *callers.Caller    // The caller used to make requests in the scenario.
	mu       sync.Mutex         // Mutex to synchronize access to the user.
}

func newUser(
	name string,
	scenario scenarios.Scenario,
	caller *callers.Caller,
) *user {
	return &user{
		name:     name,
		scenario: scenario,
		caller:   caller,
	}
}

// Run is a method that makes the user execute its assigned scenario.
// It ensures that only one action (execution of the scenario) happens at a time using a lock.
//
// Parameters:
//   - ctx: The context used for managing the lifecycle of the request.
//
// This method increments the active users gauge, executes the scenario, and then decrements the active users gauge.
func (u *user) Run(ctx context.Context) {
	// Attempt to acquire a lock for this user to prevent concurrent execution.
	if !u.mu.TryLock() {
		return // Skip if unable to lock the user (i.e., another execution is in progress).
	}
	defer u.mu.Unlock() // Ensure the lock is released after the execution is done.

	// Increment the ActiveUsersGauge metric to track the number of active users.
	metrics.ActiveUsersGauge.Inc()
	defer metrics.ActiveUsersGauge.Dec() // Decrement the metric once the user is done.

	// Execute the scenario commands for this user. If an error occurs, log it.
	if err := u.scenario.Commands(ctx, u.caller); err != nil {
		log.Printf("error with execute scenario (user: %s, scenario: %s): %v", u.name, u.scenario.Name, err)
	}
}

// Destroy is a method to destroy the user, allowing for custom logic to be added for cleanup.
// It uses a lock to ensure that no other actions can happen during the destroy process.
//
// Parameters:
//   - ctx: The context used for the destruction process.
func (u *user) Destroy(ctx context.Context) {
	// Acquire a lock to ensure the destroy process is synchronized.
	u.mu.Lock()
	defer u.mu.Unlock()

	// Implement any custom destroy logic here. Currently, there is no additional logic.
}
