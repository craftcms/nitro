package main

import (
	"log"
	"os"

	"github.com/docker/docker/client"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/add"
	"github.com/craftcms/nitro/command/apply"
	"github.com/craftcms/nitro/command/clean"
	"github.com/craftcms/nitro/command/completion"
	"github.com/craftcms/nitro/command/composer"
	"github.com/craftcms/nitro/command/context"
	"github.com/craftcms/nitro/command/create"
	"github.com/craftcms/nitro/command/database"
	"github.com/craftcms/nitro/command/destroy"
	"github.com/craftcms/nitro/command/hosts"
	"github.com/craftcms/nitro/command/initialize"
	"github.com/craftcms/nitro/command/logs"
	"github.com/craftcms/nitro/command/npm"
	"github.com/craftcms/nitro/command/queue"
	"github.com/craftcms/nitro/command/restart"
	"github.com/craftcms/nitro/command/ssh"
	"github.com/craftcms/nitro/command/start"
	"github.com/craftcms/nitro/command/stop"
	"github.com/craftcms/nitro/command/trust"
	"github.com/craftcms/nitro/command/update"
	"github.com/craftcms/nitro/command/validate"
	"github.com/craftcms/nitro/command/version"
	"github.com/craftcms/nitro/command/xon"
	nitro "github.com/craftcms/nitro/pkg/client"

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

	// create the docker client
	client, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	// get the port for the nitrod API
	apiPort := "5000"
	if os.Getenv("NITRO_API_PORT") != "" {
		apiPort = os.Getenv("NITRO_API_PORT")
	}

	// create the nitrod gRPC API
	nitrod, err := nitro.NewClient("127.0.0.1", apiPort)
	if err != nil {
		log.Fatal(err)
	}

	// create the "terminal" for capturing output
	term := terminal.New()

	// register all of the commands
	commands := []*cobra.Command{
		initialize.NewCommand(client, term),
		hosts.New(home, term),
		create.New(client, term),
		add.New(home, client, term),
		start.New(client, term),
		stop.New(client, term),
		queue.NewCommand(home, client, term),
		logs.NewCommand(home, client, term),
		ssh.New(home, client, term),
		restart.New(client, term),
		update.New(client, term),
		destroy.NewCommand(home, client, term),
		composer.NewCommand(client, term),
		npm.New(client, term),
		completion.New(),
		apply.New(home, client, nitrod, term),
		clean.NewCommand(home, client, term),
		context.New(home, client, term),
		trust.New(client, term),
		version.New(client, nitrod, term),
		validate.New(home, client, term),
		database.New(home, client, term),
		hosts.New(home, term),
		xon.NewCommand(home, client, term),
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
