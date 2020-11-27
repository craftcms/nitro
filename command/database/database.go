package database

import (
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const exampleText = `  # import a database from a backup
  nitro db import mybackup.sql`

// New returns the db commands for importing, backing up, and adding databases
func New(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "db",
		Short:   "Manage databases",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(importCommand(docker, output))

	return cmd
}
