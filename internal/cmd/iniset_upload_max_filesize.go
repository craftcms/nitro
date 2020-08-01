package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitrod"
	"github.com/craftcms/nitro/internal/client"
	"github.com/craftcms/nitro/internal/nitro"
)

var inisetUploadMaxFilesizeCommand = &cobra.Command{
	Use:   "upload_max_filesize",
	Short: "Change upload_max_filesize",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		runner := nitro.NewMultipassRunner("multipass")
		ip := nitro.IP(machine, runner)
		c := client.NewClient(ip, "50051")
		php := config.GetString("php", flagPhpVersion)

		resp, err := c.PhpIniSettings(cmd.Context(), &nitrod.ChangePhpIniSettingRequest{
			Version: php,
			Setting: nitrod.PhpIniSetting_UPLOAD_MAX_FILESIZE,
			Value:   args[0],
		})
		if err != nil {
			return err
		}

		fmt.Println(resp.Message)

		return nil
	},
}
