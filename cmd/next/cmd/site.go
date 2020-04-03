package cmd

import (
	"github.com/spf13/cobra"
)

var siteCommand = &cobra.Command{
	Use:   "site",
	Short: "Perform site commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		return nil
	},
}
