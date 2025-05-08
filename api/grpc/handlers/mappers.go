package handlers

import (
	"load-generation-system/internal/core"
	"load-generation-system/pkg/grpc/go/pb"
)

func (service *Service) mapScenarioToCore(scenario *pb.Scenario) core.ScenarioDetails {
	return core.ScenarioDetails{
		Name:        scenario.Name,
		Description: scenario.Description,
	}
}

func (service *Service) mapStartFromCore(start core.OperationStart) *pb.AttackResponse {
	return &pb.AttackResponse{
		Response: &pb.AttackResponse_Start{
			Start: &pb.OperationStart{
				Id:          start.ID,
				AttackId:    start.AttackID,
				IncrementId: start.IncrementID,
				WaitTimeSec: float32(start.WaitTimeSec), // nolint: unconvertable types from int64 to float32
				Scenarios:   start.Scenarios,
			},
		},
	}
}

func (service *Service) mapStopFromCore(stop core.OperationStop) *pb.AttackResponse {
	return &pb.AttackResponse{
		Response: &pb.AttackResponse_Stop{
			Stop: &pb.OperationStop{
				AttackId:    stop.AttackID,
				IncrementId: stop.IncrementID,
			},
		},
	}
}
