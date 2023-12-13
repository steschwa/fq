package main

import (
	"fmt"
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
		fmt.Println(err)
		os.Exit(1)
	}
}
