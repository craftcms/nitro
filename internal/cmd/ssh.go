package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

var sshCommand = &cobra.Command{
	Use:   "ssh",
	Short: "SSH into machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := config.GetString("name", flagMachineName)

		return nitro.Run(nitro.NewMultipassRunner("multipass"), nitro.SSH(name))
	},
}
