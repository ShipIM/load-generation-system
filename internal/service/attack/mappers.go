package attack

import (
	"load-generation-system/internal/core"
)

func (s *attackService) mapStartAttackToOperationStart(start core.StartAttack, attackID, incrementID int64) core.OperationStart {
	resultScenarios := make(map[string]int64)

	if start.ConstConfig != nil {
		for _, scenario := range start.ConstConfig.Scenarios {
			resultScenarios[scenario.Name] = *scenario.Counter
		}
	}

	if start.LinearConfig != nil {
		for _, scenario := range start.LinearConfig.Scenarios {
			resultScenarios[scenario.Name] += start.LinearConfig.StartCounter
		}
	}

	return core.OperationStart{
		AttackID:    attackID,
		IncrementID: incrementID,
		WaitTimeSec: start.WaitTimeSec,
		Scenarios:   resultScenarios,
	}
}
