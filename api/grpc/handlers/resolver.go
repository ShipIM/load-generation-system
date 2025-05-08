package handlers

import (
	"context"
	"load-generation-system/pkg/grpc/go/pb"
	grpc "load-generation-system/pkg/grpc/server"
)

type (
	Resolver struct {
		grpcServer *grpc.Server
		// services
		attackService *Service
	}
)

func NewResolver(
	grpcServer *grpc.Server,
	attackService *Service,
) *Resolver {
	return &Resolver{
		grpcServer: grpcServer,

		attackService: attackService,
	}
}

func (resolver Resolver) Run(ctx context.Context) {
	resolver.init()

	resolver.grpcServer.Run(ctx)
}

func (resolver Resolver) init() {
	pb.RegisterAttackServer(resolver.grpcServer.Server(), resolver.attackService)
}

func (resolver Resolver) Shutdown() error {
	resolver.grpcServer.Stop()

	return nil
}
