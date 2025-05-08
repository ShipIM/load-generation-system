package api

import (
	grpc_r "load-generation-system/api/grpc/handlers"
	manager "load-generation-system/api/rest/manager/handlers"
	"load-generation-system/pkg/rest"
)

type ManagerContainer struct {
	Server       rest.Server
	RestResolver *manager.Resolver
	GRPCResolver *grpc_r.Resolver
}

func NewManagerContainer(
	server rest.Server,
	restResolver *manager.Resolver,
	gRPCResolver *grpc_r.Resolver,
) ManagerContainer {
	return ManagerContainer{
		Server:       server,
		RestResolver: restResolver,
		GRPCResolver: gRPCResolver,
	}
}
