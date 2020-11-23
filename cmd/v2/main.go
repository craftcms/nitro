package main

import (
	"log"
	"os"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/initialize"
	"github.com/craftcms/nitro/command/update"
	"github.com/craftcms/nitro/pkg/cmd/apply"
	"github.com/craftcms/nitro/pkg/cmd/complete"
	"github.com/craftcms/nitro/pkg/cmd/composer"
	"github.com/craftcms/nitro/pkg/cmd/context"
	"github.com/craftcms/nitro/pkg/cmd/db"
	"github.com/craftcms/nitro/pkg/cmd/destroy"
	"github.com/craftcms/nitro/pkg/cmd/exec"
	"github.com/craftcms/nitro/pkg/cmd/initcmd"
	"github.com/craftcms/nitro/pkg/cmd/ls"
	"github.com/craftcms/nitro/pkg/cmd/npm"
	"github.com/craftcms/nitro/pkg/cmd/restart"
	"github.com/craftcms/nitro/pkg/cmd/start"
	"github.com/craftcms/nitro/pkg/cmd/stop"
	"github.com/craftcms/nitro/pkg/cmd/trust"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/terminal"
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
	env, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// set any global flags
	flags := rootCommand.PersistentFlags()
	flags.StringP("environment", "e", env, "The environment")

	// create the docker client
	client, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	// create the "terminal" for capturing output
	term := terminal.New()

	// register all of the commands
	commands := []*cobra.Command{
		initcmd.InitCommand,
		initialize.New(client, term),
		stop.StopCommand,
		start.StartCommand,
		destroy.DestroyCommand,
		restart.RestartCommand,
		ls.LSCommand,
		composer.ComposerCommand,
		npm.NPMCommand,
		complete.CompleteCommand,
		apply.ApplyCommand,
		context.ContextCommand,
		exec.ExecCommand,
		trust.TrustCommand,
		db.DBCommand,
		update.New(),
	}

	// add the commands
	rootCommand.AddCommand(commands...)
}

func main() {
	// execute the root command
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
