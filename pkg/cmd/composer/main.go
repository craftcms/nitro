package composer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/client"
)

// ComposerCommand is the command for creating new development environments
var ComposerCommand = &cobra.Command{
	Use:   "composer",
	Short: "Run composer actions",
	RunE:  composerMain,
	Example: `  # run composer install in a current directory
  nitro composer

  # updating a composer project outside of the current directory
  nitro composer ./project-dir --version 2 --update`,
}

func composerMain(cmd *cobra.Command, args []string) error {
	var path string
	switch len(args) {
	case 0:
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("unable to get the current directory, %w", err)
		}

		path, err = filepath.Abs(wd)
		if err != nil {
			return fmt.Errorf("unable to find the absolute path, %w", err)
		}
	default:
		var err error
		path, err = filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("unable to find the absolute path, %w", err)
		}
	}

	// create the new client
	nitro, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("unable to create a client for docker, %w", err)
	}

	// determine the default action
	action := "install"
	if cmd.Flag("update").Value.String() == "true" {
		action = "update"
	}

	return nitro.Composer(cmd.Context(), path, cmd.Flag("version").Value.String(), action)
}

func init() {
	flags := ComposerCommand.Flags()

	flags.BoolP("update", "u", false, "Run composer update instead of install")
	flags.StringP("version", "v", "1", "The composer version to use")
}
