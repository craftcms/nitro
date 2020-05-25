package cmd

import (
	"github.com/spf13/cobra"
)

var installCommand = &cobra.Command{
	Use:       "install",
	Short:     "Install software",
	ValidArgs: []string{"composer", "mailhog", "postgres"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	installCommand.AddCommand(mailhogCommand, composerCommand, postgresCommand, mysqlCommand)
}
