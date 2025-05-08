package main

import (
	"load-generation-system/cmd/manager"
	"load-generation-system/cmd/node"
	"load-generation-system/version"
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
)

// @Version 1.0.0
// @Title load-generation-system
// @Description Load Generation Service Documentation
// @ContactName Ilya Shipunov
// @Server http://localhost:8080 localhost
func main() {
	app := &cli.App{
		Usage: "LOAD GENERATION SYSTEM",
		Commands: []*cli.Command{
			&manager.Cmd,
			&node.Cmd,
		},
		Flags:   nil,
		Version: version.Version + " (" + version.GitCommit + ")",
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			return cli.ShowAppHelp(c)
		},
		Before: func(ctx *cli.Context) error {
			// place to configure some global stuff
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Println(err.Error())
	}
}
