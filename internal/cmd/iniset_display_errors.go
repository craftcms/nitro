package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/client"
	"github.com/craftcms/nitro/internal/config"
	"github.com/craftcms/nitro/internal/nitrod"
)

var inisetDisplayErrorsCommand = &cobra.Command{
	Use:   "display_errors",
	Short: "Enable or disable display_errors",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName
		c, err := client.NewDefaultClient(machine)
		if err != nil {
			return err
		}

		resp, err := c.PhpIniSettings(cmd.Context(), &nitrod.ChangePhpIniSettingRequest{
			Version: config.GetString("php", flagPhpVersion),
			Setting: nitrod.PhpIniSetting_DISPLAY_ERRORS,
			Value:   args[0],
		})
		if err != nil {
			return err
		}

		if !flagSilent {
			fmt.Println(resp.Message)
		}

		return nil
	},
}

func init() {
	inisetDisplayErrorsCommand.Flags().BoolVar(&flagSilent, "silent", false, "Run command with no output")
}
