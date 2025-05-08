package core

import "context"

// Operation represents a unit of work that can either be started, stopped, or killed.
type Operation struct {
	Start *OperationStart // Represents the operation to start an attack.
	Stop  *OperationStop  // Represents the operation to stop an attack.
	Kill  *OperationKill  // Represents the operation to kill an attack.
}

// OperationStart contains the details required to start an attack operation.
// It includes the attack ID, increment ID, wait time before starting, and the scenarios to be executed.
type OperationStart struct {
	ID          string           // Unique identifier for this operation.
	AttackID    int64            // ID of the attack to start.
	IncrementID int64            // ID of the increment to start.
	WaitTimeSec float64          // Time (in seconds) to wait before starting the operation.
	Scenarios   map[string]int64 // A map of scenario names and their respective counters.
}

// OperationStop represents the operation to stop an attack or an increment.
type OperationStop struct {
	AttackID    int64  // ID of the attack to stop.
	IncrementID *int64 // ID of the increment to stop (optional).
}

// OperationKill represents an operation to immediately kill a node.
type OperationKill struct {
	// No fields necessary for killing a node, as this operation is an immediate termination.
}

// Node represents a unit of execution that can manage and perform operations on attacks.
// Each node can start attacks, stop attacks, and acknowledge operations.
type Node interface {
	// Start initializes the node's work.
	Start(ctx context.Context)

	// StartAttack starts an attack based on the provided start details.
	StartAttack(start OperationStart) error

	// StopAttack stops an attack or increment based on the provided stop details.
	StopAttack(stop OperationStop) error

	// GetDetails retrieves the current details of the node, including scenarios and attacks.
	GetDetails() NodeDetails

	// AckOperation acknowledges the completion of an operation (start or stop).
	AckOperation()
}
