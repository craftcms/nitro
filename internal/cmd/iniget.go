package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/client"
	"github.com/craftcms/nitro/internal/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/internal/nitrod"
)

var inigetCommand = &cobra.Command{
	Use:   "iniget",
	Short: "Get PHP settings",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		ip := nitro.IP(machine, runner)
		c, err := client.NewClient(ip, "50051")
		if err != nil {
			return err
		}
		php := config.GetString("php", flagPhpVersion)

		resp, err := c.GetPhpIniSetting(cmd.Context(), &nitrod.GetPhpIniSettingRequest{
			Version: php,
			Setting: args[0],
		})
		if err != nil {
			return err
		}

		fmt.Println(resp.Message)

		return nil
	},
}
