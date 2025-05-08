package node

import (
	"github.com/urfave/cli/v2"
)

var cmdFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "server-host",
		Usage:   "server host",
		EnvVars: []string{"SERVER_HOST"},
		Value:   "localhost:8082",
	},
	&cli.StringFlag{
		Name:    "grpc-manager-host",
		Usage:   "grpc manager host",
		EnvVars: []string{"GRPC_MANAGER_HOST"},
		Value:   "localhost:5000",
	},
	&cli.IntFlag{
		Name:    "metrics-port",
		Usage:   "metrics port",
		EnvVars: []string{"METRICS_PORT"},
		Value:   4002,
	},
	&cli.StringFlag{
		Name:    "node-name",
		Usage:   "node name",
		EnvVars: []string{"NODE_NAME"},
		Value:   "node1",
	},
	&cli.Int64Flag{
		Name:    "generator-users-per-client",
		Usage:   "generator users per client",
		EnvVars: []string{"GENERATOR_USERS_PER_CLIENT"},
		Value:   15,
	},
	&cli.Int64Flag{
		Name:    "generator-min-idle-conn-timeout-sec",
		Usage:   "generator min idle conn timeout sec",
		EnvVars: []string{"GENERATOR_MIN_IDLE_CONN_TIMEOUT_SEC"},
		Value:   10,
	},
	&cli.Int64Flag{
		Name:    "generator-max-idle-conn-timeout-sec",
		Usage:   "generator max idle conn timeout sec",
		EnvVars: []string{"GENERATOR_MAX_IDLE_CONN_TIMEOUT_SEC"},
		Value:   60,
	},
}
