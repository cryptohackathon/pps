package main

import (
	"fmt"
	"github.com/ZenGo-X/fe-hackaton-demo/cli/subcommands"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			&subcommands.Keygen,
			&subcommands.SendSignal,
			&subcommands.Search,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
