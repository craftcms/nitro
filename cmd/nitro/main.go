package main

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/attach"
	"github.com/craftcms/nitro/internal/bootstrap"
	"github.com/craftcms/nitro/internal/delete"
	"github.com/craftcms/nitro/internal/destroy"
	"github.com/craftcms/nitro/internal/executor"
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

func run(args []string) {
	e := executor.NewSyscallExecutor("multipass")

	app := &cli.App{
		Name:                 "nitro",
		UsageText:            "nitro [global options] command [command options] [arguments...]",
		Usage:                "Local Craft CMS on Tap.",
		Version:              "1.0.0",
		Description:          "Nitro creates virtual machines with Multipass and provides a CLI for common DevOps tasks.",
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "machine",
				Aliases:     []string{"m"},
				Value:       "nitro-dev",
				Usage:       "Provide a machine name",
				DefaultText: "nitro-dev",
			},
		},
		Commands: []*cli.Command{
			initialize.Command(),
			bootstrap.Command(e),
			host.Command(),
			attach.Command(),
			ssh.Command(e),
			xdebug.Command(),
			start.Command(),
			stop.Command(),
			delete.Command(),
			destroy.Command(),
			password.Command(e),
			sql.Command(e),
			redis.Command(e),
			update.Command(),
			logs.Command(e),
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

	if err := app.RunContext(ctx, args); err != nil {
		log.Fatal(err)
	}
}

func main() {
	run(os.Args)
}
