package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"

	"github.com/pixelandtonic/nitro/internal/action"
	"github.com/pixelandtonic/nitro/internal/command"
	"github.com/pixelandtonic/nitro/internal/initialize"
	"github.com/pixelandtonic/nitro/internal/sql"
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
		Usage:                "Develop Craft CMS websites locally with ease",
		Version:              "1.0.0",
		Description:          "A better way to develop Craft CMS applications without Docker or Vagrant",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		Flags: []cli.Flag{machineFlag},
		Commands: []*cli.Command{
			initialize.Command(),
			command.Bootstrap(executor),
			command.AddHost(executor),
			command.SSH(executor),
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
			{
				Name:  "share",
				Usage: "Create a shareable URL to access your application",
				Action: func(c *cli.Context) error {
					return errors.New("not implemented")
				},
			},
			{
				Name:  "start",
				Usage: "Start your machine",
				Action: func(c *cli.Context) error {
					return errors.New("not implemented")
				},
			},
			command.Stop(),
			command.Delete(),
			{
				Name:        "destroy",
				Usage:       "Permanently shutdown and destroy a machine",
				Description: "By default, when deleting a machine it is soft-deleted which means it can be recovered. This command will destroy the machine making it unrecoverable.",
				Action: func(c *cli.Context) error {
					return errors.New("not implemented")
				},
			},
			{
				Name:        "password",
				Usage:       "Get the database password for the user nitro",
				Description: "Regardless of the database engine, there is one password for the non-root user. This password is unique to each machine and generated on startup.",
				Action: func(c *cli.Context) error {
					return action.DatabasePassword(c, executor)
				},
			},
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
			{
				Name:        "redis",
				Usage:       "Enter a shell for redis",
				Description: "Redis is installed by default on the machine, this command will drop you immediately into a shell inside the machine to run commands.",
				Action: func(c *cli.Context) error {
					return action.RedisCLIShell(c, executor)
				},
			},
			command.Update(),
			command.Logs(executor),
			command.IP(),
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
