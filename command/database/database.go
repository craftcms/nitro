package database

import (
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/protob"
)

const exampleText = `  # import a database from a backup
  nitro db import mybackup.sql

  # backup a database
  nitro db backup

  # add a new database
  nitro db add`

// NewCommand returns the db commands for importing, backing up, and adding databases
func NewCommand(home string, docker client.CommonAPIClient, nitrod protob.NitroClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "db",
		Short:   "Manage databases",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(
		importCommand(home, docker, nitrod, output),
		backupCommand(home, docker, output),
		addCommand(docker, output),
		sshCommand(home, docker, output),
		removeCommand(docker, output),
	)

	return cmd
}
