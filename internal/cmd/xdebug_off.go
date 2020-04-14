package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
	"github.com/craftcms/nitro/validate"
)

var xdebugOffCommand = &cobra.Command{
	Use:   "off",
	Short: "Disable Xdebug on a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)
		php := config.GetString("php", flagPhpVersion)
		if err := validate.PHPVersion(php); err != nil {
			return err
		}

		disableXdebugAction, err := nitro.DisableXdebug(name, php)
		if err != nil {
			return err
		}

		return nitro.Run(nitro.NewMultipassRunner("multipass"), []nitro.Action{*disableXdebugAction})
	},
}
