package update

import (
	"bytes"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/version"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	dockerImages = map[string]string{
		"nginx:8.0-dev":                  "docker.io/craftcms/nginx:8.0-dev",
		"nginx:7.4-dev":                  "docker.io/craftcms/nginx:7.4-dev",
		"nginx:7.3-dev":                  "docker.io/craftcms/nginx:7.3-dev",
		"nginx:7.2-dev":                  "docker.io/craftcms/nginx:7.2-dev",
		"nginx:7.1-dev":                  "docker.io/craftcms/nginx:7.1-dev",
		"nitro-proxy:" + version.Version: "docker.io/craftcms/nitro-proxy:" + version.Version,
		// TODO(jasonmccallister) finish adding builds for the 7.0 images
		//"nginx:7.0-dev":                  "docker.io/craftcms/nginx:7.0-dev",
	}
)

// New returns the update command for updating images on the local machine as well as the nitro-proxy container.
func NewCommand(docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update nitro",
		Example: `  # update nitro
  nitro update`,
		RunE: func(cmd *cobra.Command, args []string) error {
			output.Info("Updating nitro…")
			ctx := cmd.Context()

			// update all of the images
			for name, image := range dockerImages {
				output.Pending("updating", name)

				// pull the image
				rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
				if err != nil {
					output.Warning()
					output.Info("  \u2717 unable to pull image", name)
					continue
				}

				buf := &bytes.Buffer{}
				if _, err := buf.ReadFrom(rdr); err != nil {
					return fmt.Errorf("unable to read the output while pulling image, %w", err)
				}

				output.Done()
			}

			// create a filter for nitro containers
			filter := filters.NewArgs()
			filter.Add("label", labels.Nitro)

			// get a list of containers
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
			if err != nil {
				return err
			}

			// check all of the containers
			for _, c := range containers {
				// only show the site containers and proxy container
				if c.Labels[labels.Type] == "dynamodb" || c.Labels[labels.Type] == "mailhog" || c.Labels[labels.Type] == "minio" || c.Labels[labels.Type] == "redis" {
					continue
				}

				// if the images match, we are up to date
				if _, ok := dockerImages[c.Image]; ok {
					continue
				}

				// stop the container if it is running
				if c.State == "running" {
					if err := docker.ContainerStop(ctx, c.ID, nil); err != nil {
						return err
					}
				}

				// remove the container
				if err := docker.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{}); err != nil {
					return err
				}
			}

			output.Info("Images updated 👍, applying changes…")

			for _, c := range cmd.Parent().Commands() {
				// set the apply command
				if c.Use == "apply" {
					if err := c.RunE(c, args); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}

	// set the flags

	return cmd
}
