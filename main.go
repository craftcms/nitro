package main

import (
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/dev/command"
)

func main() {
	app := &cli.App{
		Name:        "dev",
		Usage:       "Quickly create new machines to develop Craft CMS websites",
		Version:     "1.0.0",
		Description: "A better way to develop PHP applications without Docker or Vagrant",
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		Commands: []*cli.Command{
			command.Initialize(),
			command.Prepare(),
			command.Build(),
			command.SSH(),
			command.Update(),
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
