package main

import (
	"log"
	"os"

	"github.com/docker/docker/client"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	nitro "github.com/craftcms/nitro/client"
	"github.com/craftcms/nitro/command/add"
	"github.com/craftcms/nitro/command/apply"
	"github.com/craftcms/nitro/command/clean"
	"github.com/craftcms/nitro/command/completion"
	"github.com/craftcms/nitro/command/composer"
	"github.com/craftcms/nitro/command/context"
	"github.com/craftcms/nitro/command/craft"
	"github.com/craftcms/nitro/command/create"
	"github.com/craftcms/nitro/command/database"
	"github.com/craftcms/nitro/command/destroy"
	"github.com/craftcms/nitro/command/disable"
	"github.com/craftcms/nitro/command/edit"
	"github.com/craftcms/nitro/command/enable"
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
	"github.com/craftcms/nitro/command/xoff"
	"github.com/craftcms/nitro/command/xon"

	"github.com/craftcms/nitro/pkg/terminal"
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
		add.NewCommand(home, client, term),
		apply.NewCommand(home, client, nitrod, term),
		clean.NewCommand(home, client, term),
		completion.New(),
		composer.NewCommand(client, term),
		context.NewCommand(home, client, term),
		craft.NewCommand(home, client, term),
		create.New(client, term),
		database.NewCommand(home, client, term),
		destroy.NewCommand(home, client, term),
		disable.NewCommand(home, client, term),
		enable.NewCommand(home, client, term),
		edit.NewCommand(home, client, term),
		hosts.New(home, term),
		initialize.NewCommand(home, client, term),
		logs.NewCommand(home, client, term),
		npm.NewCommand(client, term),
		queue.NewCommand(home, client, term),
		restart.New(client, term),
		ssh.NewCommand(home, client, term),
		start.NewCommand(client, term),
		stop.New(client, term),
		trust.New(client, term),
		update.NewCommand(client, term),
		validate.NewCommand(home, client, term),
		version.New(client, nitrod, term),
		xon.NewCommand(home, client, term),
		xoff.NewCommand(home, client, term),
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
