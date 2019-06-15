package main

import (
	"os"

	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	raven.SetDSN(sentryDSN)

	app := cli.NewApp()
	app.Flags = flags
	app.Before = func(c *cli.Context) error {
		if err := before(); err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	}
	app.Action = func(c *cli.Context) error {
		if err := action(); err != nil {
			return cli.NewExitError(err, 1)
		}
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		zap.S().Fatal(err)
	}
}
