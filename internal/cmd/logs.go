package cmd

import (
	"github.com/spf13/cobra"
)

var logsCommand = &cobra.Command{
	Use:   "logs",
	Short: "Show machine logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}
