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
	rootCmd.PersistentFlags().StringVarP(&flagMachineName, "machine", "m", "", "Name of the machine.")
	rootCmd.PersistentFlags().BoolVarP(&flagDebug, "debug", "d", false, "Bypass executing the commands.")
	rootCmd.PersistentFlags().StringVarP(&flagConfigFile, "config", "f", "", "Configuration file to use.")

	// add commands to root
	rootCmd.AddCommand(
		addCommand,
		sshCommand,
		updateCommand,
		infoCommand,
		stopCommand,
		restartCommand,
		startCommand,
		machineCommand,
		logsCommand,
		xdebugCommand,
		redisCommand,
		hostsCommand,
		contextCommand,
		selfUpdateCommand,
		applyCommand,
		removeCommand,
		destroyCommand,
		editCommand,
		importCommand,
	)
	xdebugCommand.AddCommand(xdebugOnCommand, xdebugOffCommand, xdebugConfigureCommand)
	machineCommand.AddCommand(destroyCommand, createCommand, restartCommand, startCommand, stopCommand)
	hostsCommand.AddCommand(hostsAddCommand, hostsRemoveCommand, hostsShowCommand)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func loadConfig() {
	if flagConfigFile != "" {
		viper.SetConfigFile(flagConfigFile)
	} else {
		home, _ := homedir.Dir()

		viper.AddConfigPath(home + "/" + ".nitro")
		viper.SetConfigName("nitro")
		viper.SetConfigType("yaml")
	}

	_ = viper.ReadInConfig()
}
