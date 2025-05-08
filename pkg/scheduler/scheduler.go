package scheduler

import (
	"context"
	"errors"
	"load-generation-system/internal/core"
	"sync"
)

// ErrSchedulerStopped is returned when attempting to modify a scheduler that has been stopped.
var ErrSchedulerStopped = errors.New("scheduler is stopped")

// Scheduler manages the execution of timed jobs with thread-safe operations.
// It supports creating, removing, and gracefully shutting down jobs.
type Scheduler struct {
	jobs    map[string]*job // Map of active jobs keyed by their IDs
	stopped bool            // Flag indicating if scheduler is stopped

	ctx    context.Context    // Context for cancellation
	cancel context.CancelFunc // Function to cancel all jobs
	wg     sync.WaitGroup     // WaitGroup to track running jobs

	mu sync.Mutex // Mutex for thread-safe operations
}

// New creates and returns a new Scheduler instance ready to accept jobs.
// The scheduler maintains its own context for graceful shutdown.
func New() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		jobs:   make(map[string]*job),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Scheduler) NewJob(intervalSec float64, task func()) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		return "", ErrSchedulerStopped
	}

	// Create and register new job
	j := newJob(intervalSec, task)
	s.jobs[j.ID] = j
	s.wg.Add(1)

	// Start job in separate goroutine
	go j.run(s.ctx, s.wg.Done)

	return j.ID, nil
}

// RemoveJob stops and removes a job with the given ID from the scheduler.
//
// Parameters:
//   - jobID: The ID of the job to remove
//
// Returns:
//   - error: ErrSchedulerStopped if scheduler is stopped,
//     core.ErrJobNotFound if job doesn't exist,
//     nil on success
func (s *Scheduler) RemoveJob(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		return ErrSchedulerStopped
	}

	j, exists := s.jobs[jobID]
	if !exists {
		return core.ErrJobNotFound
	}

	// Stop the job and remove from tracking
	j.stopJob()
	delete(s.jobs, jobID)

	return nil
}

// Shutdown gracefully stops all jobs and prevents new jobs from being created.
// It waits for all running jobs to complete before returning.
//
// Returns:
//   - error: nil on success, no errors currently returned
func (s *Scheduler) Shutdown() error {
	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		return nil
	}

	// Mark as stopped and cancel all jobs
	s.stopped = true
	s.cancel()

	// Get copy of job IDs to avoid holding lock during removal
	jobIDs := make([]string, 0, len(s.jobs))
	for id := range s.jobs {
		jobIDs = append(jobIDs, id)
	}
	s.mu.Unlock()

	// Remove all jobs (errors ignored as we're shutting down)
	for _, id := range jobIDs {
		_ = s.RemoveJob(id)
	}

	// Wait for all jobs to complete
	s.wg.Wait()

	return nil
}
