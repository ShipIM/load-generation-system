package scenarios

import (
	"context"
	"load-generation-system/internal/service/callers"
)

// Scenario represents a load generation scenario that contains a name, description,
// and a function that executes commands during the scenario execution.
type Scenario struct {
	// Name is the name of the scenario.
	Name string

	// Description provides a textual explanation of the scenario's purpose or functionality.
	Description string

	// Commands is a function that takes a context and a caller object, and executes a series of actions or commands.
	// It is expected to return an error if something goes wrong during the scenario execution.
	Commands func(ctx context.Context, caller *callers.Caller) error
}

// New is a constructor function that creates and returns a new Scenario instance.
// It accepts the scenario's name, description, and the function (commands) to be executed for the scenario.
func New(
	name, description string, // The name and description of the scenario.
	commands func(ctx context.Context, caller *callers.Caller) error, // The function to execute the scenario's actions.
) Scenario {
	return Scenario{
		Name:        name,
		Description: description,
		Commands:    commands,
	}
}

// AvailableScenarios is a map that holds predefined load generation scenarios.
// The key is the scenario name, and the value is the Scenario struct that contains its details and commands.
var (
	AvailableScenarios = map[string]Scenario{
		testHTTP: testHTTPScen, // Example scenario: testHTTP which is a predefined scenario.
	}
)
