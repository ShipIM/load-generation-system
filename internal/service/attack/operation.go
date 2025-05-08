package attack

import (
	"load-generation-system/internal/core"
	"load-generation-system/pkg/broadcast"
	"log"
	"time"

	"slices"

	"github.com/google/uuid"
)

// StartAttack initiates a new load test attack with the given configuration.
// It handles both constant and linear ramp attack patterns.
//
// Parameters:
//   - start: Configuration details for the new attack
//
// Returns:
//   - core.AttackDetails: Details of the created attack
//   - error: Possible errors:
//   - core.ErrScenarioNotFound if invalid scenarios specified
//   - core.ErrEmptyAttack if no valid scenarios remain after validation
//   - Errors from operation distribution
//
// The method:
// 1. Generates unique attack and increment IDs
// 2. Validates and distributes scenarios to available nodes
// 3. Starts duration or linear ramp handlers if configured
// 4. Returns the created attack details
func (s *attackService) StartAttack(start core.StartAttack) (core.AttackDetails, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	operationStart := s.mapStartAttackToOperationStart(start, s.attackSeq, 0)
	if err := s.distributeStart(operationStart); err != nil {
		return core.AttackDetails{}, err
	}

	s.attackSeq++
	s.incrementSeqs[operationStart.AttackID]++

	incrementDetails := core.IncrementDetails{
		ID:        operationStart.IncrementID,
		AttackID:  operationStart.AttackID,
		Scenarios: operationStart.Scenarios,
	}
	increments := []core.IncrementDetails{incrementDetails}

	createdAt := time.Now().UTC().Truncate(time.Second)
	attackDetails := core.AttackDetails{
		ID:           operationStart.AttackID,
		Name:         start.Name,
		WaitTimeSec:  start.WaitTimeSec,
		CreatedAt:    createdAt,
		DurationSec:  start.DurationSec,
		ConstConfig:  start.ConstConfig,
		LinearConfig: start.LinearConfig,
		Increments:   increments,
	}
	attack := attack{
		details: attackDetails,
		stopBr:  broadcast.NewBroadcaster[any](),
	}
	s.attacks[operationStart.AttackID] = attack

	// Start appropriate attack pattern handlers
	if start.DurationSec != nil {
		go s.handleDuration(attack)
	}
	if start.LinearConfig != nil {
		go s.handleLinear(attack)
	}

	return attackDetails, nil
}

// StartIncrement adds a new increment to an existing attack.
//
// Parameters:
//   - start: Configuration for the new increment
//
// Returns:
//   - core.IncrementDetails: Details of the created increment
//   - error: Possible errors:
//   - core.ErrAttackNotFound if specified attack doesn't exist
//   - Errors from operation distribution
//
// The method:
// 1. Validates the parent attack exists
// 2. Generates a new increment ID
// 3. Distributes scenarios to available nodes
// 4. Updates attack details with new increment
func (s *attackService) StartIncrement(start core.OperationStart) (core.IncrementDetails, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	attack, exists := s.attacks[start.AttackID]
	if !exists {
		return core.IncrementDetails{}, core.ErrAttackNotFound
	}

	// Set increment parameters from parent attack
	start.AttackID = attack.details.ID
	start.IncrementID = s.incrementSeqs[attack.details.ID]
	start.WaitTimeSec = attack.details.WaitTimeSec

	if err := s.distributeStart(start); err != nil {
		return core.IncrementDetails{}, err
	}

	s.incrementSeqs[attack.details.ID]++

	// Create and store increment details
	incrementDetails := core.IncrementDetails{
		ID:        start.IncrementID,
		AttackID:  start.AttackID,
		Scenarios: start.Scenarios,
	}
	attack.details.Increments = append(attack.details.Increments, incrementDetails)
	s.attacks[start.AttackID] = attack

	return incrementDetails, nil
}

// distributeStart validates scenarios and distributes them across available nodes.
//
// Parameters:
//   - start: Operation details to distribute
//
// Returns:
//   - error: Validation errors if scenarios are invalid
//
// The method:
// 1. Validates all scenarios exist in the system
// 2. Removes scenarios with zero or negative amounts
// 3. Divides the workload across nodes that support each scenario
func (s *attackService) distributeStart(start core.OperationStart) error {
	if err := s.validateScenarios(start.Scenarios); err != nil {
		return err
	}

	s.divideTasks(start)
	return nil
}

// validateScenarios checks if all specified scenarios exist in the system.
//
// Parameters:
//   - scenarios: Map of scenario names to requested amounts
//
// Returns:
//   - error: Possible errors:
//   - core.ErrScenarioNotFound if any scenario doesn't exist
//   - core.ErrEmptyAttack if no valid scenarios remain
func (s *attackService) validateScenarios(scenarios map[string]int64) error {
	uniqueScenarios := s.getScenarios()
	for scenario, amount := range scenarios {
		if _, exists := uniqueScenarios[scenario]; !exists {
			return core.ErrScenarioNotFound
		}
		if amount <= 0 {
			delete(scenarios, scenario)
		}
	}
	if len(scenarios) == 0 {
		return core.ErrEmptyAttack
	}
	return nil
}

// divideTasks distributes scenario workloads across available nodes.
//
// Parameters:
//   - start: Operation containing scenarios to distribute
//
// The method:
// 1. Creates operation structures for each node
// 2. Evenly splits scenario amounts across nodes that support them
// 3. Handles remainder distribution for uneven splits
// 4. Starts the operations on each node
func (s *attackService) divideTasks(start core.OperationStart) {
	operations := make(map[string]core.OperationStart)
	for node := range s.nodes {
		operations[node] = core.OperationStart{
			ID:          uuid.NewString(),
			AttackID:    start.AttackID,
			IncrementID: start.IncrementID,
			WaitTimeSec: start.WaitTimeSec,
			Scenarios:   make(map[string]int64),
		}
	}

	for scenario, amount := range start.Scenarios {
		// Find nodes that support this scenario
		var actualNodes []string
		for nodeName, node := range s.nodes {
			for _, scenarioDetails := range node.GetDetails().Scenarios {
				if scenarioDetails.Name == scenario {
					actualNodes = append(actualNodes, nodeName)
					break
				}
			}
		}

		// Calculate split amounts
		splitAmount := amount / int64(len(actualNodes))
		remainder := amount % int64(len(actualNodes))

		// Distribute with remainder handling
		index := 0
		for _, node := range actualNodes {
			resultAmount := splitAmount
			if int64(index) < remainder {
				resultAmount += 1
				index++
			}

			if resultAmount != 0 {
				operations[node].Scenarios[scenario] = resultAmount
			}
		}
	}

	// Start operations on each node
	for nodeName, operation := range operations {
		if len(operation.Scenarios) != 0 {
			if err := s.nodes[nodeName].StartAttack(operation); err != nil {
				log.Printf("impossible to start attack on node %s: %v", nodeName, err)
			}
		}
	}
}

// StopAttack terminates an entire attack and all its increments.
//
// Parameters:
//   - attackID: ID of the attack to stop
//
// Returns:
//   - error: Possible errors:
//   - core.ErrAttackNotFound if attack doesn't exist
func (s *attackService) StopAttack(attackID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	attack, exists := s.attacks[attackID]
	if !exists {
		return core.ErrAttackNotFound
	}

	operation := core.OperationStop{
		AttackID: attackID,
	}
	s.distributeStop(operation)

	attack.stopBr.Broadcast(nil)
	delete(s.attacks, attackID)

	return nil
}

// StopIncrement terminates a specific increment within an attack.
//
// Parameters:
//   - attackID: ID of the parent attack
//   - incrementID: ID of the increment to stop
//
// Returns:
//   - error: Possible errors:
//   - core.ErrAttackNotFound if attack doesn't exist
//   - core.ErrIncrementNotFound if increment doesn't exist
func (s *attackService) StopIncrement(attackID, incrementID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.stopIncrement(attackID, incrementID)
}

// stopIncrement is the internal implementation of increment stopping.
//
// Parameters:
//   - attackID: ID of the parent attack
//   - incrementID: ID of the increment to stop
//
// Returns:
//   - error: Possible errors (same as StopIncrement)
//
// The method:
// 1. Validates attack and increment existence
// 2. Distributes stop commands to all nodes
// 3. Updates attack details or removes attack if last increment
func (s *attackService) stopIncrement(attackID, incrementID int64) error {
	attack, exists := s.attacks[attackID]
	if !exists {
		return core.ErrAttackNotFound
	}

	// Find increment position
	increments := make(map[int64]int)
	for i, increment := range attack.details.Increments {
		increments[increment.ID] = i
	}
	i, exists := increments[incrementID]
	if !exists {
		return core.ErrIncrementNotFound
	}

	// Distribute stop command
	operation := core.OperationStop{
		AttackID:    attackID,
		IncrementID: &incrementID,
	}
	s.distributeStop(operation)

	// Update or remove attack
	if len(attack.details.Increments) > 1 {
		attack.details.Increments = slices.Delete(attack.details.Increments, i, i+1)
		s.attacks[attackID] = attack
	} else {
		delete(s.attacks, attackID)
	}

	return nil
}

// distributeStop sends stop commands to all nodes for an operation.
//
// Parameters:
//   - stop: Operation stop details to distribute
//
// The method handles logging of any node-specific stop failures.
func (s *attackService) distributeStop(stop core.OperationStop) {
	for nodeName, node := range s.nodes {
		if err := node.StopAttack(stop); err != nil {
			log.Printf("impossible to stop attack on node %s: %v", nodeName, err)
		}
	}
}
