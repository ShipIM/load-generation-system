package attack

import (
	"load-generation-system/internal/core"
	"log"
	"math"
	"time"
)

// handleDuration manages the timed execution of an attack with a fixed duration.
// It automatically stops the attack when the duration elapses or when a stop signal is received.
//
// Parameters:
//   - attack: The attack configuration containing duration and stop broadcaster
func (s *attackService) handleDuration(attack attack) {
	// Create timer with attack duration
	timer := time.NewTimer(time.Duration(*attack.details.DurationSec) * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// Duration elapsed - stop the attack
			if err := s.StopAttack(attack.details.ID); err != nil {
				log.Printf("error stopping attack %d: %v", attack.details.ID, err)
			}
			return
		case <-attack.stopBr.Subscribe():
			// Stop signal received
			return
		}
	}
}

// handleLinear manages linear ramp-up attacks with configurable parameters.
// It dynamically adjusts the load based on the linear configuration parameters.
//
// Parameters:
//   - attack: The attack configuration containing linear ramp parameters
//
// The method handles four configuration scenarios:
//  1. Both duration and step specified - computes interval automatically
//  2. Only duration specified - computes step size based on remaining time
//  3. Only step specified - uses specified step with optional interval
//  4. Neither specified - returns immediately (invalid configuration)
//
// The linear ramp continues until reaching endCounter or receiving a stop signal.
func (s *attackService) handleLinear(attack attack) {
	// Extract linear configuration parameters
	startCounter := attack.details.LinearConfig.StartCounter
	endCounter := attack.details.LinearConfig.EndCounter
	step := attack.details.LinearConfig.CounterStep
	interval := attack.details.LinearConfig.StepIntervalSec
	duration := attack.details.LinearConfig.WarmUpSec

	var computedStep int64
	computedInterval := 1.0 // Default interval if not specified
	rangeVal := endCounter - startCounter

	// Calculate step and interval based on configuration combination
	switch {
	case duration == nil && step == nil:
		// Invalid configuration - no way to determine progression
		return
	case duration != nil && step != nil:
		// Both duration and step specified - compute optimal interval
		computedInterval = float64(*duration) * float64(*step) / float64(rangeVal)
		computedStep = *step
	case duration != nil && step == nil:
		// Only duration specified - compute step based on interval
		if interval != nil {
			computedInterval = float64(*interval)
		}
		computedStep = int64(math.Ceil(float64(rangeVal) * computedInterval / float64(*duration)))
	case duration == nil && step != nil:
		// Only step specified - use specified values
		if interval != nil {
			computedInterval = float64(*interval)
		}
		computedStep = *step
	}

	currentCounter := startCounter
	var totalElapsedTime float64

	// Initialize timer with computed interval
	timer := time.NewTimer(time.Duration(computedInterval * float64(time.Second)))
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			// Check if we've reached the target load
			if currentCounter >= endCounter {
				return
			}

			// Dynamic step adjustment when only duration is specified
			if step == nil {
				computedStep = int64(math.Ceil(float64(endCounter-currentCounter) *
					computedInterval / (float64(*duration) - totalElapsedTime)))
			} else if computedStep > endCounter-currentCounter {
				// Adjust final step to exactly reach endCounter
				computedStep = endCounter - currentCounter
			}

			// Prepare scenarios with computed step size
			scenarios := make(map[string]int64)
			for _, scenario := range attack.details.LinearConfig.Scenarios {
				scenarios[scenario.Name] = computedStep
			}

			// Start new increment with calculated load
			incrementStart := core.OperationStart{
				AttackID:  attack.details.ID,
				Scenarios: scenarios,
			}
			if _, err := s.StartIncrement(incrementStart); err != nil {
				log.Printf("error starting increment: %v", err)
				return
			}

			// Update tracking variables
			totalElapsedTime += computedInterval
			currentCounter += computedStep

			// Adjust interval for final step if needed
			if duration != nil && totalElapsedTime+computedInterval > float64(*duration) {
				computedInterval = float64(*duration) - totalElapsedTime
			}

			// Reset timer for next increment
			timer.Reset(time.Duration(computedInterval * float64(time.Second)))

		case <-attack.stopBr.Subscribe():
			// Stop signal received
			return
		}
	}
}
