//go:build wireinject
// +build wireinject

package inject

import (
	"context"
	api "load-generation-system/api/node"

	"github.com/google/wire"
	"github.com/urfave/cli/v2"
)

func InitializeNode(c *cli.Context, appCtx context.Context) (api.NodeContainer, error) {
	wire.Build(
		nodeSet,
		api.NewNodeContainer,
	)
	return api.NodeContainer{}, nil
}
