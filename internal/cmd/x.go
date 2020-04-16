package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/craftcms/nitro/config"
)

var xCommand = &cobra.Command{
	Use: "x",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)
		path, err := exec.LookPath("multipass")
		if err != nil {
			return err
		}

		// load the config file
		var configFile config.Config
		if err := viper.Unmarshal(&configFile); err != nil {
			return err
		}

		for _, site := range configFile.Sites {
			c := exec.Command(path, "exec", name, "--", "sudo", "bash", "/opt/nitro/scripts/site-exists.sh", site.Hostname)
			output, err := c.Output()
			if err != nil {
				return err
			}
			if strings.Contains(string(output), "exists") {
				fmt.Println(site.Hostname, "exists")
			}
		}

		return nil
	},
}
