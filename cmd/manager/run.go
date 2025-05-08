package manager

import (
	"context"
	"load-generation-system/api/inject"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

var Cmd = cli.Command{
	Name:  "manager",
	Usage: "Run manager",
	Flags: cmdFlags,
	OnUsageError: func(c *cli.Context, err error, isSubCommand bool) error {
		return cli.ShowCommandHelp(c, "manager")
	},
	Action: run,
}

func run(c *cli.Context) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
			return
		case s := <-sig:
			log.Printf("signal %s received", s.String())
			cancel()
		}
	}()

	app, err := inject.InitializeManager(c, ctx)
	if err != nil {
		log.Fatalf("main: cannot initialize load manager: %s", err.Error())
	}

	go app.Server.Run(ctx)

	go app.GRPCResolver.Run(ctx)

	<-ctx.Done()

	_ = app.Server.Shutdown(ctx)
	_ = app.GRPCResolver.Shutdown()

	return nil
}
