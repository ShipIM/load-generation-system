package handlers

import (
	"load-generation-system/internal/core"
	"load-generation-system/internal/scenarios"
	"load-generation-system/pkg/grpc/go/pb"
)

func (gateway *attackGateway) mapStartToCore(start *pb.OperationStart) core.OperationStart {
	return core.OperationStart{
		ID:          start.Id,
		AttackID:    start.AttackId,
		IncrementID: start.IncrementId,
		WaitTimeSec: float64(start.WaitTimeSec),
		Scenarios:   start.Scenarios,
	}
}

func (gateway *attackGateway) mapStopToCore(stop *pb.OperationStop) core.OperationStop {
	return core.OperationStop{
		AttackID:    stop.AttackId,
		IncrementID: stop.IncrementId,
	}
}

func (gateway *attackGateway) mapScenario(scenario scenarios.Scenario) *pb.Scenario {
	return &pb.Scenario{
		Name:        scenario.Name,
		Description: scenario.Description,
	}
}
