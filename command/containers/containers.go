package containers

import (
	"errors"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/terminal"
)

const exampleText = `  # add custom containers
  nitro containers`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "containers",
		Short:   "Add custom containers",
		Example: exampleText,
		RunE: func(cmd *cobra.Command, args []string) error {
			// load the config
			_, err := config.Load(home)
			if errors.Is(err, config.ErrNoConfigFile) {
				// TODO(jasonmccallister) prompt the user to run init
				return err
			}

			// ask for the image
			// ask for the tag (default to latest)
			// check if the tag is available or
			// expose ports?
			// check if the port is available
			// use volumes?
			// setup a custom env file?
			// prompt for the container name

			// save the config and apply

			return nil
		},
	}

	return cmd
}
