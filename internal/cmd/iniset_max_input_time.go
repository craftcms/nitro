package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitrod"
	"github.com/craftcms/nitro/internal/client"
	"github.com/craftcms/nitro/internal/nitro"
)

var inisetMaxInputTimeCommand = &cobra.Command{
	Use:   "max_input_time",
	Short: "Change max_input_time",
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

		resp, err := c.PhpIniSettings(cmd.Context(), &nitrod.ChangePhpIniSettingRequest{
			Version: php,
			Setting: nitrod.PhpIniSetting_MAX_INPUT_TIME,
			Value:   args[0],
		})
		if err != nil {
			return err
		}

		fmt.Println(resp.Message)

		return nil
	},
}
