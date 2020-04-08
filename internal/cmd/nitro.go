package cmd

import (
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagConfigFile  string
	flagMachineName string
	flagDebug       bool

	nitroCommand = &cobra.Command{
		Use:          "nitro",
		Short:        "Local Craft CMS on tap",
		Long:         `TODO add the long description`,
		SilenceUsage: true,
	}
)

func init() {
	cobra.OnInitialize(loadConfig)

	// set persistent flags on the root command
	nitroCommand.PersistentFlags().StringVarP(&flagMachineName, "machine", "m", "", "name of machine")
	nitroCommand.PersistentFlags().BoolVarP(&flagDebug, "debug", "d", false, "bypass executing the commands")
	nitroCommand.PersistentFlags().StringVarP(&flagConfigFile, "config-file", "f", "", "configuration file to use")

	// add commands to root
	nitroCommand.AddCommand(
		siteCommand,
		sshCommand,
		redisCommand,
		updateCommand,
		xdebugCommand,
		infoCommand,
		sqlCommand,
		stopCommand,
		startCommand,
		ipCommand,
		machineCommand,
		logsCommand,
		completionCmd,
	)
	siteCommand.AddCommand(siteAddCommand, siteRemoveCommand)
	xdebugCommand.AddCommand(xdebugOnCommand, xdebugOffCommand)
	machineCommand.AddCommand(destroyCommand, createCommand, restartCommand, startCommand, stopCommand)
	logsCommand.AddCommand(logsNginxCommand, logsDatabaseCommand)
	restartDatabaseCommand.AddCommand(servicesDatabaseRestartCommand)
	restartCommand.AddCommand(restartDatabaseCommand)
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
