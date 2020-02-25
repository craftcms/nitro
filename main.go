package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"sort"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/phpdev/action"
	"github.com/pixelandtonic/phpdev/command"
)

func main() {

	executor := action.NewSyscallExecutor("multipass")

	app := &cli.App{
		Name:        "phpdev",
		Usage:       "Develop Craft CMS websites locally with ease",
		Version:     "1.0.0",
		Description: "A better way to develop Craft CMS applications without Docker or Vagrant",
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "machine",
				Aliases:     []string{"m"},
				Value:       "phpdev",
				Usage:       "Provide a machine name",
				DefaultText: "phpdev",
			},
		},
		Commands: []*cli.Command{
			command.Initialize(),
			command.SSH(executor),
			command.Update(),
			command.Install(),
			command.Delete(),
			command.Stop(),
			command.Mount(),
		},
	}

	// find the path to multipass and set value in context
	multipass, err := exec.LookPath("multipass")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.WithValue(context.Background(), "multipass", multipass)

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
