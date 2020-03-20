package app

import (
	"github.com/urfave/cli/v2"

	"github.com/craftcms/nitro/internal/command"
)

var (
	// Version is the application version that is passed at runtime.
	Version = "1.0.0"
)

func NewApp(r command.Runner) *cli.App {
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
			command.Initialize(r),
			command.Add(r),
			command.Hosts(r),
			command.Remove(r),
			command.Attach(r),
			command.SSH(r),
			command.XOn(r),
			command.XOff(r),
			command.Info(r),
			command.Start(r),
			command.Stop(r),
			command.Destroy(r),
			command.SQL(r),
			command.Redis(r),
			command.Update(r),
			command.Logs(r),
			command.IP(r),
			command.Services(r),
			command.Refresh(r),
		},
	}
}
