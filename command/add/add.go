package add

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	// ErrExample is used when we want to share an error
	ErrExample = fmt.Errorf("some example error")
)

const exampleText = `  # add the current project as a site
  nitro add

  # add a directory as the site
  nitro add my-project`

// New is used for scaffolding new commands
func New(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Short:   "Add a new site",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			output.Info("Adding site...")

			// get the environment
			site := config.Site{}

			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// get working directory or provided arg
			var dir string
			switch len(args) {
			case 1:
				dir = filepath.Join(wd, args[0])
			default:
				dir = filepath.Clean(wd)
			}

			// set the sites directory but make the path relative
			site.Path = strings.Replace(dir, home, "~", 1)

			output.Success("added site at", site.Path)

			// prompt for the php version
			versions := []string{"7.4", "7.3", "7.2", "7.1"}
			selected, err := output.Select(cmd.InOrStdin(), "Choose a PHP version: ", versions)
			if err != nil {
				return err
			}

			// set the version of php
			site.PHP = versions[selected]

			fmt.Println(site)

			// prompt for the webroot
			// prompt to enable xdebug for the site
			return nil
		},
	}

	return cmd
}
