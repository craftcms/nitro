package cmd

import (
	"github.com/spf13/cobra"
)

var (
	restartCommand = &cobra.Command{
		Use:   "restart",
		Short: "Restart a machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)
