package nitro

import (
	"log"
	"os"

	nitroclient "github.com/craftcms/nitro/client"
	"github.com/craftcms/nitro/command/add"
	"github.com/craftcms/nitro/command/alias"
	"github.com/craftcms/nitro/command/apply"
	"github.com/craftcms/nitro/command/blackfire"
	"github.com/craftcms/nitro/command/bridge"
	"github.com/craftcms/nitro/command/completion"
	"github.com/craftcms/nitro/command/composer"
	"github.com/craftcms/nitro/command/container"
	"github.com/craftcms/nitro/command/craft"
	"github.com/craftcms/nitro/command/create"
	"github.com/craftcms/nitro/command/database"
	"github.com/craftcms/nitro/command/destroy"
	"github.com/craftcms/nitro/command/disable"
	"github.com/craftcms/nitro/command/edit"
	"github.com/craftcms/nitro/command/enable"
	"github.com/craftcms/nitro/command/extensions"
	"github.com/craftcms/nitro/command/hosts"
	"github.com/craftcms/nitro/command/iniset"
	"github.com/craftcms/nitro/command/initialize"
	"github.com/craftcms/nitro/command/keys"
	"github.com/craftcms/nitro/command/logs"
	"github.com/craftcms/nitro/command/ls"
	"github.com/craftcms/nitro/command/npm"
	"github.com/craftcms/nitro/command/php"
	"github.com/craftcms/nitro/command/portcheck"
	"github.com/craftcms/nitro/command/queue"
	"github.com/craftcms/nitro/command/remove"
	"github.com/craftcms/nitro/command/restart"
	"github.com/craftcms/nitro/command/resume"
	"github.com/craftcms/nitro/command/selfupdate"
	"github.com/craftcms/nitro/command/share"
	"github.com/craftcms/nitro/command/ssh"
	"github.com/craftcms/nitro/command/start"
	"github.com/craftcms/nitro/command/stop"
	"github.com/craftcms/nitro/command/suspend"
	"github.com/craftcms/nitro/command/trust"
	"github.com/craftcms/nitro/command/update"
	"github.com/craftcms/nitro/command/version"
	"github.com/craftcms/nitro/command/xoff"
	"github.com/craftcms/nitro/command/xon"
	"github.com/craftcms/nitro/pkg/downloader"
	"github.com/craftcms/nitro/pkg/flags"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/docker/docker/client"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "nitro",
	Short: "Speedy local dev environment for Craft CMS.",
	Long: `Nitro is a console-based tool that manages Docker for local PHP development.

Version: ` + version.Version,
	RunE:         rootMain,
	SilenceUsage: true,
	Version:      version.Version,
}

func rootMain(command *cobra.Command, _ []string) error {
	return command.Help()
}

func NewCommand() *cobra.Command {
	// get the users home directory
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	// create the docker client
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}

	// get the port for the nitrod API
	apiPort := "5000"
	if os.Getenv("NITRO_API_PORT") != "" {
		apiPort = os.Getenv("NITRO_API_PORT")
	}

	// create the nitrod gRPC API
	nitrod, err := nitroclient.NewClient("127.0.0.1", apiPort)
	if err != nil {
		log.Fatal(err)
	}

	// create the "terminal" for capturing output
	term := terminal.New()

	// create the downloaded for creating projects
	downloader := downloader.NewDownloader()

	// register the commands
	commands := []*cobra.Command{
		add.NewCommand(home, docker, term),
		alias.NewCommand(home, docker, term),
		apply.NewCommand(home, docker, nitrod, term),
		blackfire.NewCommand(home, docker, term),
		bridge.NewCommand(home, docker, term),
		completion.NewCommand(),
		composer.NewCommand(home, docker, term),
		container.NewCommand(home, docker, term),
		craft.NewCommand(home, docker, term),
		create.NewCommand(home, docker, downloader, term),
		database.NewCommand(home, docker, nitrod, term),
		destroy.NewCommand(home, docker, term),
		disable.NewCommand(home, docker, term),
		enable.NewCommand(home, docker, term),
		suspend.NewCommand(home, docker, term),
		resume.NewCommand(home, docker, term),
		edit.NewCommand(home, docker, term),
		extensions.NewCommand(home, docker, term),
		hosts.NewCommand(home, term),
		iniset.NewCommand(home, docker, term),
		initialize.NewCommand(home, docker, term),
		keys.NewCommand(home, docker, term),
		logs.NewCommand(home, docker, term),
		ls.NewCommand(home, docker, term),
		npm.NewCommand(home, docker, term),
		php.NewCommand(home, docker, term),
		portcheck.NewCommand(term),
		queue.NewCommand(home, docker, term),
		remove.NewCommand(home, docker, term),
		restart.NewCommand(home, docker, term),
		selfupdate.NewCommand(term),
		share.NewCommand(home, docker, term),
		ssh.NewCommand(home, docker, term),
		start.NewCommand(home, docker, term),
		stop.NewCommand(home, docker, term),
		trust.NewCommand(home, docker, term),
		update.NewCommand(home, docker, term),
		version.NewCommand(home, docker, nitrod, term),
		xon.NewCommand(home, docker, term),
		xoff.NewCommand(home, docker, term),
	}

	// add the commands
	rootCommand.AddCommand(commands...)

	// add the global app flag
	rootCommand.PersistentFlags().StringVarP(&flags.AppName, "app", "a", "", "the app to use for the command")

	return rootCommand
}
