package core

import (
	"time"
)

// StartAttack represents the configuration for starting a new attack.
// It includes details such as the attack name, wait time, duration, and configurations for different types of attack strategies.
type StartAttack struct {
	Name         string        // Name of the attack.
	WaitTimeSec  float64       // Time to wait between attack executions (in seconds).
	DurationSec  *int64        // Duration of the attack (in seconds). If nil, no duration limit.
	ConstConfig  *ConstConfig  // Configuration for constant attack strategy.
	LinearConfig *LinearConfig // Configuration for linear attack strategy.
}

// ConstConfig defines the configuration for a constant attack, where scenarios are run at a fixed rate.
type ConstConfig struct {
	Scenarios []Scenario // List of scenarios to run during the attack.
}

// LinearConfig defines the configuration for a linear attack, where the attack gradually increases over time.
type LinearConfig struct {
	WarmUpSec       *int64     // Time to warm up before the attack starts (in seconds).
	StartCounter    int64      // The starting counter value for the attack.
	EndCounter      int64      // The ending counter value for the attack.
	CounterStep     *int64     // The step increment/decrement of the counter.
	StepIntervalSec *int64     // Interval between each step (in seconds).
	Scenarios       []Scenario // List of scenarios to run during the attack.
}

// ScenarioDetails contains metadata about a scenario, including its name and description.
type ScenarioDetails struct {
	Name        string // Name of the scenario.
	Description string // Description of the scenario.
}

// IncrementDetails provides details about an increment in the attack, such as the increment ID and associated scenarios.
type IncrementDetails struct {
	ID        int64            // Unique ID for the increment.
	AttackID  int64            // ID of the attack that this increment belongs to.
	Scenarios map[string]int64 // A map of scenarios with their respective counters.
}

// AttackDetails contains all the details about an attack, including the configuration and its increments.
type AttackDetails struct {
	ID           int64              // Unique ID of the attack.
	Name         string             // Name of the attack.
	WaitTimeSec  float64            // Wait time before starting the attack.
	CreatedAt    time.Time          // Time when the attack was created.
	DurationSec  *int64             // Duration for the attack (in seconds).
	ConstConfig  *ConstConfig       // Constant attack configuration.
	LinearConfig *LinearConfig      // Linear attack configuration.
	Increments   []IncrementDetails // List of increments associated with the attack.
}

// NodeDetails contains details about a node, including its name, whether it's active, and the scenarios it can run.
type NodeDetails struct {
	Name      string            // Name of the node.
	IsActive  bool              // Indicates whether the node is active.
	Scenarios []ScenarioDetails // List of scenarios available for the node.
	Attacks   []AttackDetails   // List of attacks assigned to the node.
}

// Scenario represents an individual scenario that can be executed during an attack.
type Scenario struct {
	Name    string // Name of the scenario.
	Counter *int64 // A counter that tracks the number of times the scenario has been executed.
}

// AttackService defines the operations available for managing and controlling attacks. It includes methods for
// starting and stopping attacks, increments, and retrieving attack details and node information.
type AttackService interface {
	// StartAttack starts a new attack based on the provided configuration.
	StartAttack(start StartAttack) (AttackDetails, error)

	// StartIncrement starts a new increment for the given operation start configuration.
	StartIncrement(start OperationStart) (IncrementDetails, error)

	// StopAttack stops the attack with the specified ID.
	StopAttack(attackID int64) error

	// StopIncrement stops the increment for the specified attack and increment IDs.
	StopIncrement(attackID, incrementID int64) error

	// GetAttacks retrieves a list of all the current attacks.
	GetAttacks() []AttackDetails

	// GetScenarios retrieves a list of all available scenarios.
	GetScenarios() []ScenarioDetails

	// ListNodes retrieves a list of all nodes in the system.
	ListNodes() []NodeDetails

	// AddNode adds a new node to the system.
	AddNode(node Node) error

	// RemoveNode removes an existing node from the system.
	RemoveNode(node Node) error
}
