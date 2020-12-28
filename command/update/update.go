package update

import (
	"bytes"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/terminal"
)

// New returns the update command for updating images on the
// local machine
func New(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update Docker images",
		Example: `  # update docker images
  nitro update`,
		RunE: func(cmd *cobra.Command, args []string) error {
			output.Info("Updating...")

			images := []string{"docker.io/craftcms/nginx:8.0-dev", "docker.io/craftcms/nginx:7.4-dev", "docker.io/craftcms/nginx:7.3-dev", "docker.io/craftcms/nginx:7.2-dev", "docker.io/craftcms/nginx:7.1-dev", "docker.io/craftcms/nginx:7.0-dev"}

			for _, image := range images {
				output.Pending("updating", image)

				// pull the image
				rdr, err := docker.ImagePull(cmd.Context(), image, types.ImagePullOptions{All: false})
				if err != nil {
					output.Warning()
					output.Info("  \u2717 unable to pull image", image)
					continue
				}

				buf := &bytes.Buffer{}
				if _, err := buf.ReadFrom(rdr); err != nil {
					return fmt.Errorf("unable to read the output from pulling the image, %w", err)
				}

				output.Done()
			}

			output.Info("Images updated 👍")

			return nil
		},
	}

	// set the flags

	return cmd
}
