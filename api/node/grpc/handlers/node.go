package handlers

import (
	"context"
	"io"
	"load-generation-system/internal/core"
	"load-generation-system/internal/scenarios"
	"load-generation-system/internal/service/generator"
	"load-generation-system/pkg/grpc/go/pb"
	"log"
)

// attackGateway implements the core.AttackGateway interface and manages the communication
// between the node and the central attack service via gRPC streaming.
// It handles operation commands (start/stop/kill) and maintains the load generation state.
type attackGateway struct {
	attackClient  pb.AttackClient    // gRPC client for attack service communication
	loadGenerator core.LoadGenerator // Load generator implementation for executing attacks

	nodeName string // Identifier for this node

	sendCh chan *pb.AttackRequest  // Channel for outgoing messages to the attack service
	recvCh chan *pb.AttackResponse // Channel for incoming messages from the attack service
	doneCh chan any                // Channel to signal when the gateway has stopped
}

func NewGateway(
	attackClient pb.AttackClient,
	config generator.Config,
	nodeName string,
) core.AttackGateway {
	return &attackGateway{
		attackClient:  attackClient,
		loadGenerator: generator.New(config),
		nodeName:      nodeName,
		sendCh:        make(chan *pb.AttackRequest),
		recvCh:        make(chan *pb.AttackResponse),
		doneCh:        make(chan any),
	}
}

// Start initiates the attack gateway operation by:
// 1. Establishing a gRPC stream with the attack service
// 2. Starting goroutines for message handling
// 3. Sending initial handshake with node information
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//
// Returns:
//   - chan any: Channel that will be closed when the gateway stops
func (g *attackGateway) Start(ctx context.Context) chan any {
	// Establish gRPC stream with attack service
	stream, err := g.attackClient.StreamAttack(ctx)
	if err != nil {
		log.Printf("failed to start gRPC attack stream: %v", err)
	}

	// Start goroutines for stream handling
	go g.runSender(ctx, stream)
	go g.runReceiver(ctx, stream)
	go g.runHandler(ctx)

	// Prepare available scenarios information
	scenariosInfo := make([]*pb.Scenario, 0, len(scenarios.AvailableScenarios))
	for _, scenario := range scenarios.AvailableScenarios {
		scenariosInfo = append(scenariosInfo, g.mapScenario(scenario))
	}

	// Send initial handshake to identify this node
	g.sendCh <- &pb.AttackRequest{
		Request: &pb.AttackRequest_Handshake{
			Handshake: &pb.Handshake{
				NodeName:  g.nodeName,
				Scenarios: scenariosInfo,
			},
		},
	}

	return g.doneCh
}

// runHandler processes incoming messages from the attack service and coordinates
// the appropriate response actions. It runs in a dedicated goroutine.
//
// Parameters:
//   - ctx: Context for cancellation
func (g *attackGateway) runHandler(ctx context.Context) {
	// Ensure cleanup on exit
	defer func() {
		g.loadGenerator.Stop() // Stop any ongoing load generation
		close(g.doneCh)        // Signal gateway shutdown
	}()

	for {
		select {
		case <-ctx.Done():
			log.Printf("app context done, stopping handler")
			return
		case resp, ok := <-g.recvCh:
			if !ok {
				log.Printf("recv channel closed, stopping handler")
				return
			}

			// Handle each response in a separate goroutine to avoid blocking
			go func(resp *pb.AttackResponse) {
				var request *pb.AttackRequest

				switch val := resp.Response.(type) {
				case *pb.AttackResponse_Start:
					request = g.handleStart(val.Start)
				case *pb.AttackResponse_Stop:
					request = g.handleStop(val.Stop)
				case *pb.AttackResponse_Kill:
					// Nil request signals the sender to close the stream
					g.sendCh <- nil
					return
				}

				g.sendCh <- request
			}(resp)
		}
	}
}

// handleStart processes a start operation command from the attack service.
//
// Parameters:
//   - start: The start operation details
//
// Returns:
//   - *pb.AttackRequest: Acknowledgment to send back to the service
func (g *attackGateway) handleStart(start *pb.OperationStart) *pb.AttackRequest {
	attackStart := g.mapStartToCore(start)

	if err := g.loadGenerator.StartAttack(attackStart); err != nil {
		log.Printf("failed to start attack: %v", err)
	}

	return &pb.AttackRequest{
		Request: &pb.AttackRequest_Acknowledge{
			Acknowledge: &pb.Acknowledge{},
		},
	}
}

// handleStop processes a stop operation command from the attack service.
//
// Parameters:
//   - stop: The stop operation details
//
// Returns:
//   - *pb.AttackRequest: Acknowledgment to send back to the service
func (g *attackGateway) handleStop(stop *pb.OperationStop) *pb.AttackRequest {
	attackStop := g.mapStopToCore(stop)
	if err := g.loadGenerator.StopAttack(attackStop); err != nil {
		log.Printf("failed to stop attack: %v", err)
	}

	return &pb.AttackRequest{
		Request: &pb.AttackRequest_Acknowledge{
			Acknowledge: &pb.Acknowledge{},
		},
	}
}

// runSender manages outgoing messages to the attack service via the gRPC stream.
// It runs in a dedicated goroutine.
//
// Parameters:
//   - ctx: Context for cancellation
//   - stream: The gRPC stream client
func (g *attackGateway) runSender(
	ctx context.Context,
	stream pb.Attack_StreamAttackClient,
) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("app context done, stopping sender")
			return
		case <-stream.Context().Done():
			log.Printf("stream context done, stopping sender")
			return
		case req := <-g.sendCh:
			if req == nil {
				// Nil request signals to close the stream
				if err := stream.CloseSend(); err != nil {
					log.Printf("error closing stream: %v", err)
				}
				log.Printf("send channel closed, stopping sender")
				return
			}

			if err := stream.Send(req); err != nil {
				log.Printf("error sending request: %v", err)
			}
		}
	}
}

// runReceiver manages incoming messages from the attack service via the gRPC stream.
// It runs in a dedicated goroutine and closes the recvCh channel when done.
//
// Parameters:
//   - ctx: Context for cancellation
//   - stream: The gRPC stream client
func (g *attackGateway) runReceiver(
	ctx context.Context,
	stream pb.Attack_StreamAttackClient,
) {
	defer close(g.recvCh) // Ensure channel is closed when receiver stops

	for {
		select {
		case <-ctx.Done():
			log.Printf("app context done, stopping receiver")
			return
		case <-stream.Context().Done():
			log.Printf("stream context done, stopping receiver")
			return
		default:
		}

		resp, err := stream.Recv()
		if err == io.EOF {
			log.Printf("stream ended by server, stopping receiver")
			return
		}
		if err != nil {
			log.Printf("error receiving response: %v", err)
			continue
		}

		g.recvCh <- resp
	}
}
