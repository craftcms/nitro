package main

import (
	"os"

	"github.com/craftcms/nitro/pkg/cmd/destroy"
	"github.com/craftcms/nitro/pkg/cmd/initcmd"
	"github.com/craftcms/nitro/pkg/cmd/restart"
	"github.com/craftcms/nitro/pkg/cmd/start"
	"github.com/craftcms/nitro/pkg/cmd/stop"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:          "nitro",
	Short:        "Local Craft CMS dev made easy",
	Long:         `Nitro is a command-line tool focused on making local Craft CMS development quick and easy.`,
	RunE:         rootMain,
	SilenceUsage: true,
}

func rootMain(command *cobra.Command, _ []string) error {
	return command.Help()
}

func init() {
	flags := rootCommand.PersistentFlags()

	// set global flags
	flags.String("environment", "nitro-dev", "The name of the environment")

	// register all of the commands
	commands := []*cobra.Command{
		initcmd.InitCommand,
		stop.StopCommand,
		start.StartCommand,
		destroy.DestroyCommand,
		restart.RestartCommand,
	}

	rootCommand.AddCommand(commands...)
}

func main() {
	// Execute the root command.
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
