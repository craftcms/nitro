package cmd

import (
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:          "nitro",
	Short:        "Local Craft CMS dev made easy",
	Long:         `Nitro is a command-line tool focused on making local Craft development quick and easy.`,
	SilenceUsage: true,
}

func init() {
	cobra.OnInitialize(loadConfig)

	// set persistent flags on the root command
	rootCmd.PersistentFlags().StringVarP(&flagMachineName, "machine", "m", "", "Name of a machine.")
	rootCmd.PersistentFlags().BoolVarP(&flagDebug, "debug", "d", false, "Show command output and do not execute.")

	// add commands to root
	rootCmd.AddCommand(
		initCommand,
		addCommand,
		sshCommand,
		keysCommand,
		updateCommand,
		infoCommand,
		stopCommand,
		restartCommand,
		startCommand,
		logsCommand,
		xdebugCommand,
		redisCommand,
		contextCommand,
		selfUpdateCommand,
		applyCommand,
		removeCommand,
		destroyCommand,
		editCommand,
		hostsCommand,
		renameCommand,
		dbCommand,
		completionCmd,
	)
	xdebugCommand.AddCommand(xdebugOnCommand, xdebugOffCommand, xdebugConfigureCommand)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func loadConfig() {
	home, _ := homedir.Dir()

	viper.AddConfigPath(home + "/" + ".nitro")
	viper.SetConfigType("yaml")

	defaultMachine := os.Getenv("NITRO_DEFAULT_MACHINE")

	if flagMachineName != "" {
		viper.SetConfigName(flagMachineName)
	} else if defaultMachine != "" {
		flagMachineName = defaultMachine
		viper.SetConfigName(defaultMachine)
	} else {
		flagMachineName = "nitro-dev"
		viper.SetConfigName("nitro-dev")
	}

	_ = viper.ReadInConfig()
}
