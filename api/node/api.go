package api

import (
	"load-generation-system/api/node/rest/handlers"
	"load-generation-system/internal/core"
	"load-generation-system/pkg/rest"

	"google.golang.org/grpc"
)

type ManagerConn struct {
	Conn *grpc.ClientConn
}

type GRPCConn struct {
	Conn *grpc.ClientConn
}

type NodeContainer struct {
	AttackGateway     core.AttackGateway
	Server            rest.Server
	RestResolver      *handlers.Resolver
	ManagerConnection ManagerConn
	GRPCConnection    GRPCConn
}

func NewNodeContainer(
	attackGateway core.AttackGateway,
	server rest.Server,
	resolver *handlers.Resolver,
	managerConnection ManagerConn,
	grpcConnection GRPCConn,
) NodeContainer {
	return NodeContainer{
		AttackGateway:     attackGateway,
		Server:            server,
		RestResolver:      resolver,
		ManagerConnection: managerConnection,
		GRPCConnection:    grpcConnection,
	}
}
