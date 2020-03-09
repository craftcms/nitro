package app

import (
	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/attach"
	"github.com/craftcms/nitro/internal/bootstrap"
	"github.com/craftcms/nitro/internal/command"
	"github.com/craftcms/nitro/internal/destroy"
	"github.com/craftcms/nitro/internal/executor"
	"github.com/craftcms/nitro/internal/host"
	"github.com/craftcms/nitro/internal/initialize"
	"github.com/craftcms/nitro/internal/ip"
	"github.com/craftcms/nitro/internal/logs"
	"github.com/craftcms/nitro/internal/redis"
	"github.com/craftcms/nitro/internal/sql"
	"github.com/craftcms/nitro/internal/ssh"
	"github.com/craftcms/nitro/internal/start"
	"github.com/craftcms/nitro/internal/stop"
	"github.com/craftcms/nitro/internal/update"
	"github.com/craftcms/nitro/internal/xdebug"
)

var (
	// Version is the application version that is passed at runtime.
	Version = "1.0.0"
)

func NewApp(e executor.Executor) *cli.App {
	return &cli.App{
		Name:        "nitro",
		UsageText:   "nitro [global options] command [command options] [arguments...]",
		Usage:       "Local Craft CMS on Tap.",
		Version:     Version,
		Description: "Nitro creates virtual machines with Multipass and provides a CLI for common DevOps tasks.",
		Action:      cli.ShowAppHelp,
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
			bootstrap.Command(),
			host.Command(),
			command.Remove(),
			attach.Command(),
			ssh.Command(e),
			xdebug.CommandOn(e),
			xdebug.CommandOff(e),
			start.Command(),
			stop.Command(),
			destroy.Command(),
			sql.Command(e),
			redis.Command(e),
			update.Command(),
			logs.Command(e),
			ip.Command(),
		},
	}
}
