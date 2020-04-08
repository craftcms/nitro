package cmd

import (
	"github.com/spf13/cobra"
)

var siteCommand = &cobra.Command{
	Use:   "site",
	Short: "Perform site commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}
