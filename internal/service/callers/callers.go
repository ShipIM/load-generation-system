package callers

import (
	"load-generation-system/internal/core"
	"load-generation-system/internal/service/callers/test"
)

// Caller is the client for making requests to the services. It manages all
// communication with the target services endpoints during load generation.
//
// Fields:
//   - TestCaller: The actual implementation that calls the Test service endpoints
//   - State: Current state of the user
type Caller struct {
	TestCaller test.TestCaller // Implementation for calling Test service
	State      core.State      // Current user state
}

// NewCaller creates a new client instance for calling target services.
//
// Parameters:
//   - httpClient: Configured HTTP client for communicating
//
// Returns:
//   - *Caller: Initialized client ready to call target services endpoints
func NewCaller(httpClient core.Client) *Caller {
	return &Caller{
		TestCaller: test.NewCaller(httpClient),
	}
}
