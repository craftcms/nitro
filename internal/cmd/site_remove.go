package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var siteRemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove a site from machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("machine", flagMachineName)

		return nitro.Run(nitro.NewMultipassRunner("multipass"), nitro.Empty(name))
	},
}
