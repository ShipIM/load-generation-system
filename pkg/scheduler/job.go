package scheduler

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// job represents a scheduled task that runs either at fixed intervals or continuously.
// It provides mechanisms to start, stop, and manage the execution of the task.
type job struct {
	ID       string        // Unique identifier for the job
	interval time.Duration // Interval between executions (0 means run continuously)
	task     func()        // The function to execute
	stop     chan any      // Channel to signal job termination
}

func newJob(intervalSec float64, task func()) *job {
	return &job{
		ID:       uuid.NewString(), // Generate unique ID for the job
		interval: time.Duration(intervalSec * float64(time.Second)),
		task:     task,
		stop:     make(chan any), // Buffered channel for stop signals
	}
}

// run starts the job's execution loop. It runs in a goroutine and will either:
// - Execute continuously if interval is 0
// - Execute at fixed intervals if interval > 0
// The provided done function is called when the job stops to notify the scheduler.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - done: Callback function to notify when job completes (typically scheduler cleanup)
func (j *job) run(ctx context.Context, done func()) {
	// Ensure done callback is called when we exit
	defer done()

	// Continuous execution mode (interval = 0)
	if j.interval == 0 {
		for {
			select {
			case <-j.stop: // Explicit stop signal
				return
			case <-ctx.Done(): // Context cancellation
				return
			default:
				j.task() // Execute task immediately
				// Note: No delay between executions in continuous mode
			}
		}
	}

	// Timed execution mode (interval > 0)
	timer := time.NewTimer(0) // Initial timer fires immediately
	defer timer.Stop()        // Ensure timer resources are cleaned up

	for {
		select {
		case <-timer.C:
			// Execute task in a goroutine to avoid blocking the timer
			go j.task()

			// Calculate next execution time aligned to the interval
			now := time.Now()
			next := now.Truncate(j.interval).Add(j.interval)
			timer.Reset(next.Sub(now)) // Reset timer for next interval

		case <-j.stop: // Explicit stop signal
			return
		case <-ctx.Done(): // Context cancellation
			return
		}
	}
}

// stopJob signals the job to stop execution by closing the stop channel.
// This is safe to call multiple times as channel closing is idempotent.
func (j *job) stopJob() {
	close(j.stop) // Signal the run loop to exit
}
