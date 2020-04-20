package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var contextCommand = &cobra.Command{
	Use:   "context",
	Short: "View the current Nitro context",
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile := viper.ConfigFileUsed()
		if configFile == "" {
			return errors.New("no configuration file being used")
		}

		// open the file
		data, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}

		fmt.Println("Using config file: ", configFile)
		fmt.Println("------")
		fmt.Print(string(data))
		fmt.Println("------")
		return nil
	},
}
