package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/nitro/internal/action"
	"github.com/pixelandtonic/nitro/internal/bootstrap"
	"github.com/pixelandtonic/nitro/internal/command"
	"github.com/pixelandtonic/nitro/internal/host"
	"github.com/pixelandtonic/nitro/internal/initialize"
	"github.com/pixelandtonic/nitro/internal/ip"
	"github.com/pixelandtonic/nitro/internal/logs"
	"github.com/pixelandtonic/nitro/internal/password"
	"github.com/pixelandtonic/nitro/internal/redis"
	"github.com/pixelandtonic/nitro/internal/sql"
	"github.com/pixelandtonic/nitro/internal/ssh"
	"github.com/pixelandtonic/nitro/internal/start"
	"github.com/pixelandtonic/nitro/internal/stop"
	"github.com/pixelandtonic/nitro/internal/update"
	"github.com/pixelandtonic/nitro/internal/x"
)

var machineFlag = &cli.StringFlag{
	Name:        "machine",
	Aliases:     []string{"m"},
	Value:       "nitro-dev",
	Usage:       "Provide a machine name",
	DefaultText: "nitro-dev",
}

func main() {
	executor := action.NewSyscallExecutor("multipass")

	app := &cli.App{
		Name:                 "nitro",
		UsageText:            "nitro [global options] command [command options] [arguments...]",
		Usage:                "Develop Craft CMS applications locally with ease",
		Version:              "1.0.0",
		Description:          "An easier way to develop Craft CMS applications without Docker or Vagrant",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		Flags: []cli.Flag{machineFlag},
		Commands: []*cli.Command{
			initialize.Command(),
			bootstrap.Command(executor),
			host.Command(executor),
			ssh.Command(executor),
			{
				Name:        "xdebug",
				Usage:       "Enable of disable xdebug on the machine",
				Description: "Calling xdebug will default to enabling the extension, if the flag --disable is passed it will be disabled",
				Action: func(c *cli.Context) error {
					return errors.New("not implemented")
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "disable",
						Usage:       "Disable xdebug",
						Value:       false,
						DefaultText: "false",
					},
				},
			},
			start.Command(),
			stop.Command(),
			command.Delete(),
			{
				Name:        "destroy",
				Usage:       "Permanently shutdown and destroy a machine",
				Description: "By default, when deleting a machine it is soft-deleted which means it can be recovered. This command will destroy the machine making it unrecoverable.",
				Action: func(c *cli.Context) error {
					return errors.New("not implemented")
				},
			},
			password.Command(executor),
			{
				Name:        "php",
				Usage:       "Install a specific version of PHP",
				Description: "The bootstrap command defaults to the latest version of PHP, this command allows you to install alternative versions of PHP.",
				Action: func(c *cli.Context) error {
					if c.Args().Len() > 0 {
						return errors.New("you must provide a version of PHP")
					}

					// TODO validate the version of PHP

					return errors.New("not implemented")
				},
				Hidden: true,
			},
			sql.Command(executor),
			redis.Command(executor),
			update.Command(),
			logs.Command(executor),
			ip.Command(),
			{
				Name: "multiple",
				Action: func(c *cli.Context) error {
					return x.MultipleCommands(c)
				},
				Hidden: true,
			},
		},
	}

	// find the path to multipass and set value in context
	multipass, err := exec.LookPath("multipass")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.WithValue(context.Background(), "multipass", multipass)

	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
