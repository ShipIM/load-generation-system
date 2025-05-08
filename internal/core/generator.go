package core

// State represents the configuration or state of the system at a particular point in time.
// It contains a map of parameters (`Params`) that can hold various dynamic configuration values.
// The parameters can be of any type (`any`), providing flexibility for different types of configurations.
type State struct {
	Params map[string]any // A map containing configuration parameters as key-value pairs.
}

// LoadGenerator is an interface that defines the operations needed to manage the lifecycle of a load generation process.
// It includes methods to start and stop attacks, as well as a method to stop the generator itself.
type LoadGenerator interface {
	// StartAttack begins the attack process for a given `OperationStart` configuration.
	// It triggers the load generation based on the provided details.
	//
	// Parameters:
	//   - start: An `OperationStart` struct that contains details about the attack to start.
	//
	// Returns:
	//   - An error if the attack could not be started; otherwise, nil.
	StartAttack(start OperationStart) error

	// StopAttack halts the attack process for a given `OperationStop` configuration.
	// It stops a running attack based on the provided stop details.
	//
	// Parameters:
	//   - stop: An `OperationStop` struct that contains details about the attack to stop.
	//
	// Returns:
	//   - An error if the attack could not be stopped; otherwise, nil.
	StopAttack(stop OperationStop) error

	// Stop terminates the load generator itself, stopping any ongoing operations.
	// This is typically used to gracefully shut down the load generation process.
	Stop()
}
