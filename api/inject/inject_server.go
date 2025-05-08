package inject

import (
	"load-generation-system/api/grpc/handlers"
	restHandlers "load-generation-system/api/rest/manager/handlers"
	"load-generation-system/internal/core"
	grpcserver "load-generation-system/pkg/grpc/server"
	"load-generation-system/pkg/rest"

	"load-generation-system/internal/service/attack"

	"github.com/google/wire"
	"github.com/urfave/cli/v2"
)

// wire set for loading the server.
var managerSet = wire.NewSet( // nolint
	provideManagerServerConfig,
	restHandlers.NewResolver,
	provideManagerGRPCConfig,
	provideManagerService,
	handlers.NewResolver,
	grpcserver.New,
	provideAttackService,
	rest.New,
)

func provideManagerServerConfig(c *cli.Context) rest.Config {
	cfg := rest.Config{
		Host:        c.String("server-host"),
		MetricsPort: c.Int("metrics-port"),
	}

	return cfg
}

func provideManagerGRPCConfig(c *cli.Context) grpcserver.Config {
	return grpcserver.Config{
		Host: c.String("grpc-manager-host"),
	}
}

func provideAttackService(c *cli.Context) core.AttackService {
	return attack.NewService(
		c.Int64("recovery-interval-sec"),
	)
}

func provideManagerService(c *cli.Context, attackService core.AttackService) *handlers.Service {
	return handlers.NewService(
		attackService,
		c.Int64("op-queue-capacity"),
		c.Int64("retry-interval-sec"),
	)
}
