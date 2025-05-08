package core

import "context"

// AttackGateway is an interface that defines the operations related to the start of an attack
// through the system. It manages the lifecycle of an attack by handling the communication
// between the client and the server during the attack execution.
type AttackGateway interface {
	// Start begins the attack process, returning a channel that can be used to track the
	// completion or any other events related to the attack's lifecycle.
	//
	// Parameters:
	//   - ctx: The context that controls the lifecycle of the attack.
	//
	// Returns:
	//   - A channel (`chan any`) which can be used to signal when the attack process has finished or handle other events.
	Start(ctx context.Context) chan any
}
