package attack

import (
	"load-generation-system/internal/core"
	"log"
	"time"
)

// AddNode registers a new node with the attack service and redistributes any
// existing operations that were assigned to a previous node with the same name.
//
// Parameters:
//   - node: The node implementation to register
//
// Returns:
//   - error: Possible errors:
//   - core.ErrNodeAlreadyExists if node name is already registered
//   - Errors from failed operation redistribution
//
// The method handles node recovery scenarios by:
// 1. Checking for existing operations from previous nodes with same name
// 2. Attempting to restart those operations on the new node
// 3. Cleaning up any failed operation attempts
func (s *attackService) AddNode(node core.Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	nodeDetails := node.GetDetails()

	// Check for existing node with same name
	_, isRemoved := s.removingCancels[nodeDetails.Name]
	if _, exists := s.nodes[nodeDetails.Name]; exists && !isRemoved {
		return core.ErrNodeAlreadyExists
	}

	// Retrieve and restart any existing operations
	operations := s.retrieveOperations(nodeDetails.Name)
	for _, operation := range operations {
		if err := node.StartAttack(operation); err != nil {
			log.Printf("impossible to start attack on node %s: %v", nodeDetails.Name, err)

			// Clean up failed operation
			if err := s.stopIncrement(operation.AttackID, operation.IncrementID); err != nil {
				log.Printf("impossible to stop increment: %v", err)
			}
		}
	}

	s.nodes[nodeDetails.Name] = node
	return nil
}

// retrieveOperations gathers all operations assigned to a node and prepares them
// for redistribution. This is used both for node recovery and removal scenarios.
//
// Parameters:
//   - nodeName: Name of the node to retrieve operations from
//
// Returns:
//   - []core.OperationStart: Slice of operations that were assigned to the node
//
// Side Effects:
//   - Removes the node from the active nodes map
//   - Closes and removes the node's removal cancellation channel
func (s *attackService) retrieveOperations(nodeName string) []core.OperationStart {
	retrieved := nodeName
	var operations []core.OperationStart

	if _, exists := s.nodes[retrieved]; exists {
		// Collect all operations from the node
		nodeDetails := s.nodes[retrieved].GetDetails()
		for _, attack := range nodeDetails.Attacks {
			attackDetails := s.attacks[attack.ID].details
			for _, increment := range attack.Increments {
				operations = append(operations, core.OperationStart{
					AttackID:    attack.ID,
					IncrementID: increment.ID,
					WaitTimeSec: attackDetails.WaitTimeSec,
					Scenarios:   increment.Scenarios,
				})
			}
		}

		// Clean up node tracking
		close(s.removingCancels[retrieved])
		delete(s.nodes, retrieved)
	}

	return operations
}

// RemoveNode initiates graceful removal of a node from the service. The removal
// follows a recovery-oriented process that:
// 1. Starts a timer during which the node may reconnect
// 2. If timer expires, redistributes the node's operations
// 3. Cleans up node resources
//
// Parameters:
//   - node: The node to remove
//
// Returns:
//   - error: Currently always returns nil (may be extended in future)
func (s *attackService) RemoveNode(node core.Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	nodeDetails := node.GetDetails()
	s.startRemovingTimer(nodeDetails.Name)
	return nil
}

// startRemovingTimer initiates the node removal process with a recovery window.
// During the recovery interval:
// - The node may reconnect and cancel the removal
// - If interval elapses, operations are redistributed
//
// Parameters:
//   - nodeName: Name of the node being removed
//
// The method manages its own locking for the timer operations to avoid
// holding the main lock during the waiting period.
func (s *attackService) startRemovingTimer(nodeName string) {
	cancel := make(chan any)
	s.removingCancels[nodeName] = cancel

	timer := time.NewTimer(s.recoveryInterval)

	go func() {
		defer func() {
			timer.Stop()

			s.mu.Lock()
			delete(s.removingCancels, nodeName)
			s.mu.Unlock()
		}()

		select {
		case <-cancel:
			// Removal was canceled (node reconnected)
			return
		case <-timer.C:
			// Recovery period elapsed - redistribute operations
			s.mu.Lock()
			operations := s.retrieveOperations(nodeName)

			for _, operation := range operations {
				if err := s.distributeStart(operation); err != nil {
					log.Printf("impossible to redistribute load: %v", err)

					if err := s.stopIncrement(operation.AttackID, operation.IncrementID); err != nil {
						log.Printf("impossible to stop increment: %v", err)
					}
				}
			}
			s.mu.Unlock()
		}
	}()
}
