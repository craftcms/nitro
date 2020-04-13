package cmd

import (
	"github.com/spf13/cobra"
)

var hostsCommand = &cobra.Command{
	Use:   "hosts",
	Short: "Manage your hosts file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}
