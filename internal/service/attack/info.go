package attack

import (
	"load-generation-system/internal/core"
)

// GetAttacks retrieves details of all currently active attacks in the system.
// The returned slice is a copy of the internal attack state for thread safety.
//
// Returns:
//   - []core.AttackDetails: A slice containing details of all active attacks
//     (empty slice if no attacks are active)
func (s *attackService) GetAttacks() []core.AttackDetails {
	s.mu.RLock()
	defer s.mu.RUnlock()

	attacks := make([]core.AttackDetails, 0, len(s.attacks))
	for _, attack := range s.attacks {
		attacks = append(attacks, attack.details)
	}

	return attacks
}

// GetScenarios retrieves all unique scenario definitions available across all nodes.
// The returned slice is deduplicated by scenario name and represents the union
// of all scenarios known to registered nodes.
//
// Returns:
//   - []core.ScenarioDetails: A slice of unique scenario definitions
//     (empty slice if no nodes are registered)
func (s *attackService) GetScenarios() []core.ScenarioDetails {
	s.mu.RLock()
	defer s.mu.RUnlock()

	uniqueScenarios := s.getScenarios()

	scenarios := make([]core.ScenarioDetails, 0, len(uniqueScenarios))
	for _, scenario := range uniqueScenarios {
		scenarios = append(scenarios, scenario)
	}

	return scenarios
}

// getScenarios is an internal helper method that collects all unique scenarios
// from registered nodes into a map keyed by scenario name.
//
// Returns:
//   - map[string]core.ScenarioDetails: Map of unique scenarios by name
func (s *attackService) getScenarios() map[string]core.ScenarioDetails {
	uniqueScenarios := make(map[string]core.ScenarioDetails)
	for _, node := range s.nodes {
		nodeDetails := node.GetDetails()

		for _, scenario := range nodeDetails.Scenarios {
			uniqueScenarios[scenario.Name] = scenario
		}
	}

	return uniqueScenarios
}

// ListNodes retrieves comprehensive details about all registered nodes including:
// - Node metadata
// - Assigned attacks
// - Attack increments
// The returned data represents a consistent snapshot of the system state.
//
// Returns:
//   - []core.NodeDetails: A slice containing full details of all registered nodes
//     (empty slice if no nodes are registered)
func (s *attackService) ListNodes() []core.NodeDetails {
	s.mu.RLock()
	defer s.mu.RUnlock()

	nodes := make([]core.NodeDetails, 0, len(s.nodes))

	for _, node := range s.nodes {
		nodeDetails := node.GetDetails()

		// Enrich node details with attack information
		attacks := make([]core.AttackDetails, 0, len(nodeDetails.Attacks))
		for _, attack := range nodeDetails.Attacks {
			attackDetails := s.attacks[attack.ID].details
			attackDetails.Increments = attack.Increments
			attacks = append(attacks, attackDetails)
		}

		nodeDetails.Attacks = attacks
		nodes = append(nodes, nodeDetails)
	}

	return nodes
}
