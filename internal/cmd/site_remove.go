package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
)

var siteRemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove a site from machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = config.GetString("machine", flagMachineName)
		return errors.New("not implemented")
		// return nitro.Run(nitro.NewMultipassRunner("multipass"), nitro.Empty(name))
	},
}
