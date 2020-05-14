package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/scripts"
)

var testCommand = &cobra.Command{
	Use:   "test",
	Short: "Testing",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		mp, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		script := scripts.New(mp, machine)

		// check if the site it available
		output, err := script.Run(fmt.Sprintf(scripts.FmtDockerContainerExists, "postgres_12_5432"))
		if err != nil {
			return err
		}

		fmt.Println(output)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCommand)
}
