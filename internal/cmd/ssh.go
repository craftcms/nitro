package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/internal/nitro"
)

var sshCommand = &cobra.Command{
	Use:   "ssh",
	Short: "SSH into a machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		machine := flagMachineName

		sshAction, err := nitro.SSH(machine)
		if err != nil {
			return err
		}

		return nitro.Run(nitro.NewMultipassRunner("multipass"), []nitro.Action{*sshAction})
	},
}
