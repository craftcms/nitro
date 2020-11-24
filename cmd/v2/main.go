package main

import (
	"log"
	"os"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/apply"
	"github.com/craftcms/nitro/command/completion"
	"github.com/craftcms/nitro/command/composer"
	"github.com/craftcms/nitro/command/context"
	"github.com/craftcms/nitro/command/destroy"
	"github.com/craftcms/nitro/command/initialize"
	"github.com/craftcms/nitro/command/node"
	"github.com/craftcms/nitro/command/restart"
	"github.com/craftcms/nitro/command/start"
	"github.com/craftcms/nitro/command/stop"
	"github.com/craftcms/nitro/command/trust"
	"github.com/craftcms/nitro/command/update"
	"github.com/craftcms/nitro/command/version"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/terminal"
)

var rootCommand = &cobra.Command{
	Use:   "nitro",
	Short: "Local Craft CMS dev made easy",
	Long: `Nitro is a command-line tool focused on making local Craft CMS development quick and easy.

Version: ` + version.Version,
	RunE:         rootMain,
	SilenceUsage: true,
	Version:      version.Version,
}

func rootMain(command *cobra.Command, _ []string) error {
	return command.Help()
}

func init() {
	env, _, err := config.Load()
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
		initialize.New(client, term),
		start.New(client, term),
		stop.New(client, term),
		restart.New(client, term),
		update.New(client, term),
		destroy.New(client, term),
		composer.New(client, term),
		node.New(client, term),
		completion.New(),
		apply.New(client, term),
		context.New(client, term),
		trust.New(client, term),
		//db.DBCommand,
		//php.New(client, term),
	}

	// add the commands
	rootCommand.AddCommand(commands...)
}

func main() {
	// cobra.OnInitialize()

	// execute the root command
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
