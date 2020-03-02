package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/bootstrap"
	"github.com/craftcms/nitro/internal/delete"
	"github.com/craftcms/nitro/internal/destroy"
	executor2 "github.com/craftcms/nitro/internal/executor"
	"github.com/craftcms/nitro/internal/host"
	"github.com/craftcms/nitro/internal/initialize"
	"github.com/craftcms/nitro/internal/ip"
	"github.com/craftcms/nitro/internal/logs"
	"github.com/craftcms/nitro/internal/password"
	"github.com/craftcms/nitro/internal/redis"
	"github.com/craftcms/nitro/internal/sql"
	"github.com/craftcms/nitro/internal/ssh"
	"github.com/craftcms/nitro/internal/start"
	"github.com/craftcms/nitro/internal/stop"
	"github.com/craftcms/nitro/internal/update"
	"github.com/craftcms/nitro/internal/x"
	"github.com/craftcms/nitro/internal/xdebug"
)

var machineFlag = &cli.StringFlag{
	Name:        "machine",
	Aliases:     []string{"m"},
	Value:       "nitro-dev",
	Usage:       "Provide a machine name",
	DefaultText: "nitro-dev",
}

func main() {
	executor := executor2.NewSyscallExecutor("multipass")

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
			xdebug.Command(),
			start.Command(),
			stop.Command(),
			delete.Command(),
			destroy.Command(),
			password.Command(executor),
			sql.Command(executor),
			redis.Command(executor),
			update.Command(),
			logs.Command(executor),
			ip.Command(),
			// this command is experimental, probably not needed
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
