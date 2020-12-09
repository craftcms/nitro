package main

import (
	"errors"
	"log"
	"os"

	"github.com/docker/docker/client"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/apply"
	"github.com/craftcms/nitro/command/completion"
	"github.com/craftcms/nitro/command/composer"
	"github.com/craftcms/nitro/command/context"
	"github.com/craftcms/nitro/command/create"
	"github.com/craftcms/nitro/command/database"
	"github.com/craftcms/nitro/command/destroy"
	"github.com/craftcms/nitro/command/hosts"
	"github.com/craftcms/nitro/command/initialize"
	"github.com/craftcms/nitro/command/node"
	"github.com/craftcms/nitro/command/queue"
	"github.com/craftcms/nitro/command/restart"
	"github.com/craftcms/nitro/command/start"
	"github.com/craftcms/nitro/command/stop"
	"github.com/craftcms/nitro/command/trust"
	"github.com/craftcms/nitro/command/update"
	"github.com/craftcms/nitro/command/validate"
	"github.com/craftcms/nitro/command/version"
	"github.com/craftcms/nitro/config"
	nitro "github.com/craftcms/nitro/pkg/client"
	"github.com/craftcms/nitro/setup"

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
	// set any global flags
	flags := rootCommand.PersistentFlags()

	// set a default environment if there is a variable
	// set the default environment name
	env := "nitro-dev"
	if os.Getenv("NITRO_DEFAULT_ENVIRONMENT") != "" {
		env = os.Getenv("NITRO_DEFAULT_ENVIRONMENT")
	}

	flags.StringP("environment", "e", env, "The environment")

	// get the users home directory
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	// check for or create the config
	if _, err := config.Load(home, env); err != nil {
		if errors.Is(err, config.ErrNoConfigFile) {
			// get the home directory
			home, err := homedir.Dir()
			if err != nil {
				log.Fatal("unable to get the home directory, %w", err)
			}

			// setup the config file
			if err := setup.FirstTime(home, env); err != nil {
				log.Fatal(err)
			}
		}
	}

	// create the docker client
	client, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	// create the nitrod gRPC API
	nitrod, err := nitro.NewClient("127.0.0.1", "5000")
	if err != nil {
		log.Fatal(err)
	}

	// create the "terminal" for capturing output
	term := terminal.New()

	// register all of the commands
	commands := []*cobra.Command{
		initialize.New(client, term),
		hosts.New(home, term),
		create.New(client, term),
		start.New(client, term),
		stop.New(client, term),
		queue.New(home, client, term),
		restart.New(client, term),
		update.New(client, term),
		destroy.New(client, term),
		composer.New(client, term),
		node.New(client, term),
		completion.New(),
		apply.New(home, client, nitrod, term),
		context.New(home, client, term),
		trust.New(client, term),
		version.New(client, nitrod, term),
		validate.New(home, client, term),
		database.New(home, client, term),
		hosts.New(home, term),
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
