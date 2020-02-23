package main

import (
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/dev/command"
)

// Runner is an interface that handles
// running commands, the real use case
// is to call syscall.Exec
type Runner interface {
	Exec(argv0 string, argv []string, envv []string) (err error)
}

func main() {
	app := &cli.App{
		Name:        "dev",
		Usage:       "Quickly create new machines to develop PHP applications and websites",
		Version:     "1.0.0",
		Description: "A better way to develop PHP applications without Docker or Vagrant",
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "machine",
				Aliases:     []string{"m"},
				Value:       "dev",
				Usage:       "Provide a machine name",
				DefaultText: "dev",
			},
		},
		Commands: []*cli.Command{
			command.Initialize(),
			command.SSH(),
			command.Update(),
			command.Install(),
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
