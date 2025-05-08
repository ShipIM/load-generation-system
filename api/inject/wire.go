//go:build wireinject
// +build wireinject

package inject

import (
	"context"
	"load-generation-system/api"

	"github.com/google/wire"
	"github.com/urfave/cli/v2"
)

func InitializeManager(c *cli.Context, appCtx context.Context) (api.ManagerContainer, error) {
	wire.Build(
		managerSet,
		api.NewManagerContainer,
	)
	return api.ManagerContainer{}, nil
}
