package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/dev/command"
)

var (
	machineFlag = cli.StringFlag{
		Name:        "machine, m",
		Value:       "dev",
		Usage:       "Provide a machine name",
		DefaultText: "dev",
	}
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
			command.Update(),
			command.SSH(),
			command.Prepare(),
			command.Build(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
