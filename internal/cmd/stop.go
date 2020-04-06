package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var stopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("machine", flagMachineName)

		return nitro.Run(nitro.NewMultipassRunner("multipass"), nitro.Stop(name))
	},
}
