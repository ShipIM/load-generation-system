package node

import (
	"context"
	"load-generation-system/api/node/inject"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

var Cmd = cli.Command{
	Name:  "node",
	Usage: "Run node",
	Flags: cmdFlags,
	OnUsageError: func(c *cli.Context, err error, isSubCommand bool) error {
		return cli.ShowCommandHelp(c, "node")
	},
	Action: run,
}

func run(c *cli.Context) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	localCtx, localCancel := context.WithCancel(appCtx)
	defer localCancel()

	go func() {
		select {
		case <-appCtx.Done():
			return
		case s := <-sig:
			log.Printf("signal %s received", s.String())
			localCancel()
		}
	}()

	app, err := inject.InitializeNode(c, appCtx)
	if err != nil {
		log.Fatalf("main: cannot initialize node: %s", err.Error())
	}
	defer app.ManagerConnection.Conn.Close()
	defer app.GRPCConnection.Conn.Close()

	go app.Server.Run(appCtx)

	doneCh := app.AttackGateway.Start(localCtx)

	<-doneCh

	_ = app.Server.Shutdown(appCtx)

	return nil
}
