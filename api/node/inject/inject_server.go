package inject

import (
	"crypto/tls"
	api "load-generation-system/api/node"
	"load-generation-system/api/node/grpc/handlers"
	restHandlers "load-generation-system/api/node/rest/handlers"
	"load-generation-system/internal/core"
	"load-generation-system/internal/metrics/interceptors"
	"load-generation-system/internal/service/generator"
	pb_manager "load-generation-system/pkg/grpc/go/pb"
	"load-generation-system/pkg/rest"

	"github.com/google/wire"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// wire set for loading the node.
var nodeSet = wire.NewSet( // nolint
	provideManagerConnection,
	provideManagerClient,
	provideNodeGRPCConnection,
	provideNodeServerConfig,
	restHandlers.NewResolver,
	provideAttackGateway,
	provideGeneratorConfig,
	rest.New,
)

func provideNodeServerConfig(c *cli.Context) rest.Config {
	cfg := rest.Config{
		Host:        c.String("server-host"),
		MetricsPort: c.Int("metrics-port"),
	}

	return cfg
}

func provideManagerConnection(c *cli.Context) (api.ManagerConn, error) {
	conn, err := grpc.NewClient(
		c.String("grpc-manager-host"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return api.ManagerConn{}, err
	}
	return api.ManagerConn{Conn: conn}, nil
}

func provideManagerClient(conn api.ManagerConn) pb_manager.AttackClient {
	return pb_manager.NewAttackClient(conn.Conn)
}

func provideNodeGRPCConnection(c *cli.Context) (api.GRPCConn, error) {
	conn, err := grpc.NewClient(
		c.String("grpc-host"),
		grpc.WithUnaryInterceptor(interceptors.GRPCInterceptor),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true, // nolint
		})),
	)
	if err != nil {
		return api.GRPCConn{}, err
	}
	return api.GRPCConn{Conn: conn}, nil
}

func provideAttackGateway(
	c *cli.Context,
	attackClient pb_manager.AttackClient,
	config generator.Config,
) core.AttackGateway {
	return handlers.NewGateway(
		attackClient,
		config,
		c.String("node-name"),
	)
}

func provideGeneratorConfig(c *cli.Context) generator.Config {
	return generator.Config{
		UsersPerClient:        c.Int64("generator-users-per-client"),
		MinIdleConnTimeoutSec: c.Int64("generator-min-idle-conn-timeout-sec"),
		MaxIdleConnTimeoutSec: c.Int64("generator-max-idle-conn-timeout-sec"),
	}
}
