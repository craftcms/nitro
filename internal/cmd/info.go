package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var infoCommand = &cobra.Command{
	Use:   "info",
	Short: "Show machine info",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("machine", flagMachineName)

		return nitro.Run(nitro.NewMultipassRunner("multipass"), nitro.Info(name))
	},
}
