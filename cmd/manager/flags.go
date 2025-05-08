package manager

import (
	"github.com/urfave/cli/v2"
)

var cmdFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "server-host",
		Usage:   "server host",
		EnvVars: []string{"SERVER_HOST"},
		Value:   "localhost:8080",
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
		Value:   4000,
	},
	&cli.Int64Flag{
		Name:    "retry-interval-sec",
		Usage:   "retry interval sec",
		EnvVars: []string{"RETRY_INTERVAL_SEC"},
		Value:   10,
	},
	&cli.Int64Flag{
		Name:    "recovery-interval-sec",
		Usage:   "recovery interval sec",
		EnvVars: []string{"RECOVERY_INTERVAL_SEC"},
		Value:   60,
	},
	&cli.Int64Flag{
		Name:    "op-queue-capacity",
		Usage:   "op queue capacity",
		EnvVars: []string{"OP_QUEUE_CAPACITY"},
		Value:   1000,
	},
}
