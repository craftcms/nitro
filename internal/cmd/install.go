package cmd

import (
	"github.com/spf13/cobra"
)

// TODO if there is an argument, assume its an apt request (e.g. apt install golang)
var installCommand = &cobra.Command{
	Use:       "install",
	Short:     "Install software",
	ValidArgs: []string{"composer", "mailhog", "postgres", "mysql"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	installCommand.AddCommand(mailhogCommand, composerCommand, postgresCommand, mysqlCommand)
}
