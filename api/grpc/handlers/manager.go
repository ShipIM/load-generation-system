package handlers

import (
	"context"
	"fmt"
	"io"
	"load-generation-system/internal/core"
	"load-generation-system/internal/service/node"
	"load-generation-system/pkg/grpc/go/pb"
	"log"
)

// Service implements the gRPC AttackServer interface and handles communication
// between the load generation system and nodes.
type Service struct {
	attackService core.AttackService // Core service for attack management

	nodeOpQueueCapacity  int64 // Maximum capacity for node operation queues
	nodeRetryIntervalSec int64 // Interval in seconds between operation retries

	pb.UnimplementedAttackServer
}

func NewService(attackService core.AttackService, nodeOpQueueCapacity, nodeRetryIntervalSec int64) *Service {
	return &Service{
		attackService:        attackService,
		nodeOpQueueCapacity:  nodeOpQueueCapacity,
		nodeRetryIntervalSec: nodeRetryIntervalSec,
	}
}

// StreamAttack establishes a bidirectional stream with a node for attack coordination.
// It handles the initial handshake, node registration, and continuous operation exchange.
//
// Parameters:
//   - stream: The gRPC server stream for bidirectional communication
//
// Returns:
//   - error: Any error that occurs during stream processing
func (service *Service) StreamAttack(stream pb.Attack_StreamAttackServer) error {
	ctx := stream.Context()

	// Receive initial handshake message from node
	resp, err := stream.Recv()
	if err == io.EOF {
		return nil // Graceful stream closure
	}
	if err != nil {
		log.Printf("error receiving initial request: %v", err)
		return err
	}

	// Validate and extract handshake data
	val, ok := resp.Request.(*pb.AttackRequest_Handshake)
	if !ok {
		err := fmt.Errorf("unable to cast request to handshake")
		log.Printf("connection cannot be established: %v", err)
		return err
	}

	// Convert protobuf scenarios to core scenarios
	scenarios := make(map[string]core.ScenarioDetails)
	for _, scenario := range val.Handshake.Scenarios {
		scenarios[scenario.Name] = service.mapScenarioToCore(scenario)
	}

	// Create operation channel for this node
	ops := make(chan core.Operation, 1)

	// Create and start new node instance
	n := node.New(
		val.Handshake.NodeName,
		scenarios,
		ops,
		service.nodeOpQueueCapacity,
		service.nodeRetryIntervalSec,
	)
	n.Start(ctx)

	// Register node with the attack service
	if err := service.attackService.AddNode(n); err != nil {
		log.Printf("error creating new node: %v", err)
	}
	// Ensure node is removed when stream ends
	defer func() {
		if err := service.attackService.RemoveNode(n); err != nil {
			log.Printf("impossible to remove node: %v", err)
		}
	}()

	// Start receiver goroutine to handle acknowledgments from node
	go service.runReceiver(ctx, stream, n)

	// Start sender in main goroutine to handle operations to node
	return service.runSender(ctx, stream, ops)
}

// runReceiver handles incoming messages from the node, primarily acknowledgments
// of completed operations.
//
// Parameters:
//   - ctx: Context for cancellation
//   - stream: The gRPC server stream
//   - n: The node instance associated with this stream
func (service *Service) runReceiver(
	ctx context.Context,
	stream pb.Attack_StreamAttackServer,
	n core.Node,
) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("stream context done, stopping receiver")
			return
		default:
		}

		// Receive message from node
		resp, err := stream.Recv()
		if err == io.EOF {
			return // Graceful stream closure
		}
		if err != nil {
			log.Printf("error receiving request: %v", err)
			return
		}

		// Validate message is an acknowledgment
		if _, ok := resp.Request.(*pb.AttackRequest_Acknowledge); !ok {
			log.Printf("unable to cast request to acknowledge")
			return
		}

		// Notify node of received acknowledgment
		n.AckOperation()
	}
}

// runSender sends operations to the node through the stream.
//
// Parameters:
//   - ctx: Context for cancellation
//   - stream: The gRPC server stream
//   - ops: Channel of operations to be sent to the node
//
// Returns:
//   - error: Any error that occurs during sending
func (service *Service) runSender(
	ctx context.Context,
	stream pb.Attack_StreamAttackServer,
	ops chan core.Operation,
) error {
	for {
		select {
		case <-ctx.Done():
			log.Printf("stream context done, stopping sender")
			return nil
		case op, ok := <-ops:
			if !ok {
				log.Printf("ops channel closed, stopping sender")
				return nil
			}

			// Convert operation to appropriate protobuf response
			var response *pb.AttackResponse
			if op.Start != nil {
				response = service.mapStartFromCore(*op.Start)
			} else if op.Stop != nil {
				response = service.mapStopFromCore(*op.Stop)
			} else if op.Kill != nil {
				response = &pb.AttackResponse{
					Response: &pb.AttackResponse_Kill{
						Kill: &pb.OperationKill{},
					},
				}
			}

			// Send operation to node
			if err := stream.Send(response); err != nil {
				log.Printf("error sending response: %v", err)
				return err
			}
		}
	}
}
