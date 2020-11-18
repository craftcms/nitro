package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
	"github.com/craftcms/nitro/pkg/cmd/update"
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
	home, _ := homedir.Dir()

	viper.AddConfigPath(fmt.Sprintf("%s%c%s", home, os.PathSeparator, ".nitro"))
	viper.SetConfigType("yaml")

	// set the default machine name
	def := "nitro-dev"
	if os.Getenv("NITRO_DEFAULT_MACHINE") != "" {
		def = os.Getenv("NITRO_DEFAULT_MACHINE")
	}

	// set the config file
	viper.SetConfigName(def)

	// read the config
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(fmt.Errorf("error loading the config file, %w", err))
	}

	// set any global flags
	flags := rootCommand.PersistentFlags()
	flags.StringP("environment", "e", def, "The environment")

	// register all of the commands
	commands := []*cobra.Command{
		initcmd.InitCommand,
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
		update.UpdateCommand,
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
