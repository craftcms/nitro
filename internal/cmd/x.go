package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
)

var xCommand = &cobra.Command{
	Use:    "x",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgFile := viper.ConfigFileUsed()
		if cfgFile == "" {
			return errors.New("unable to find the config file")
		}

		filePath, err := filepath.Abs(cfgFile)
		if err != nil {
			return err
		}

		nitro, err := exec.LookPath("nitro")
		if err != nil {
			return err
		}

		fmt.Println("ok, modifying the hosts file to add sites for", config.GetString("name", ""), "(you will be prompted for your password)... ")
		hostsCmd := exec.Command("sudo", nitro, "-f", filePath, "hosts", "add")
		hostsCmd.Stdout = os.Stdout
		hostsCmd.Stderr = os.Stderr

		return hostsCmd.Run()
	},
}

func init() {
	rootCmd.AddCommand(xCommand)
}
