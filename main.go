package main

import (
	"log/slog"
	"os"

	"github.com/steschwa/fq/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "fq",
		Usage: "firebase query tool",
		Commands: []*cli.Command{
			cmd.QueryCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
