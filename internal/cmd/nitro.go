package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagMachineName string
	flagDyRun       bool

	nitroCommand = &cobra.Command{
		Use:   "nitro",
		Short: "Local Craft CMS on tap",
		Long:  `TODO add the long description`,
	}
)

func init() {
	cobra.OnInitialize(loadConfig)

	// set persistent flags on the root command
	nitroCommand.PersistentFlags().StringVarP(&flagMachineName, "machine", "m", "", "name of machine")
	nitroCommand.PersistentFlags().BoolVarP(&flagDyRun, "dry-run", "d", false, "bypass executing the commands")

	// add commands to root
	nitroCommand.AddCommand(siteCommand, sshCmd, initCommand, redisCommand, updateCommand, destroyCommand, xdebugCommand)
	siteCommand.AddCommand(siteAddCommand, siteRemoveCommand)
	xdebugCommand.AddCommand(xdebugOnCommand, xdebugOffCommand)
}

func Execute() {
	if err := nitroCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadConfig() {
	home, _ := homedir.Dir()

	viper.AddConfigPath(home + "/" + ".nitro")
	viper.SetConfigName("nitro")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
