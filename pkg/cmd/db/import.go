package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/craftcms/nitro/pkg/client"
	"github.com/spf13/cobra"
)

// ImportCommand is the command for creating new development environments
var ImportCommand = &cobra.Command{
	Use:   "import",
	Short: "Import a database",
	Args:  cobra.MinimumNArgs(1),
	RunE:  importMain,
	Example: `  # import a sql file into a database
  nitro db import filename.sql

  # use a relative path
  nitro db import ~/Desktop/backup.sql

  # use an absolute path
  nitro db import /Users/oli/Desktop/backup.sql`,
}

func importMain(cmd *cobra.Command, args []string) error {
	// create the new client
	nitro, err := client.NewClient()
	if err != nil {
		return err
	}

	path := filepath.Clean(args[0])

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening the file, %w", err)
	}

	return nitro.Import(cmd.Context(), "59876a79434e", "/app/web/readme.md", f)
}
