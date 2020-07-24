package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/api"
	"github.com/craftcms/nitro/internal/client"
	"github.com/craftcms/nitro/internal/nitro"
)

var inisetMaxfileuploadsCommand = &cobra.Command{
	Use:   "max_file_uploads",
	Short: "Change max_file_uploads",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		ip := nitro.IP(machine, runner)
		c := client.NewClient(ip, "50051")
		php := config.GetString("php", flagPhpVersion)

		resp, err := c.PhpIniSettings(cmd.Context(), &api.ChangePhpIniSettingRequest{
			Version: php,
			Setting: api.PhpIniSetting_MAX_FILE_UPLOADS,
			Value:   args[0],
		})
		if err != nil {
			return err
		}

		fmt.Println(resp.Message)

		return nil
	},
}
