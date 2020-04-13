package cmd

import (
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version = "development"

	nitroCommand = &cobra.Command{
		Use:          "nitro",
		Short:        "Local Craft CMS dev made easy",
		Long:         `nitro is a command-line tool focused on making local Craft development quick and easy.`,
		Version:      Version,
		SilenceUsage: true,
	}
)

func init() {
	cobra.OnInitialize(loadConfig)

	// set persistent flags on the root command
	nitroCommand.PersistentFlags().StringVarP(&flagMachineName, "machine", "m", "", "name of machine")
	nitroCommand.PersistentFlags().BoolVarP(&flagDebug, "debug", "d", false, "bypass executing the commands")
	nitroCommand.PersistentFlags().StringVarP(&flagConfigFile, "config", "f", "", "configuration file to use")

	// add commands to root
	nitroCommand.AddCommand(
		siteCommand,
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
		versionCommand,
	)
	siteCommand.AddCommand(siteAddCommand, siteRemoveCommand)
	xdebugCommand.AddCommand(xdebugOnCommand, xdebugOffCommand, xdebugConfigureCommand)
	machineCommand.AddCommand(destroyCommand, createCommand, restartCommand, startCommand, stopCommand)
	logsCommand.AddCommand(logsNginxCommand, logsDockerCommand, logsDatabaseCommand)
	hostsCommand.AddCommand(hostsAddCommand, hostsRemoveCommand, hostsShowCommand)
}

func Execute() {
	if err := nitroCommand.Execute(); err != nil {
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
