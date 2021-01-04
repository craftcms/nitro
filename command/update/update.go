package update

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/version"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	DockerImages = map[string]string{
		// "docker.io/craftcms/nginx:8.0-dev": "nginx:8.0-dev",
		"docker.io/craftcms/nginx:7.4-dev": "nginx:7.4-dev",
		"docker.io/craftcms/nginx:7.3-dev": "nginx:7.3-dev",
		"docker.io/craftcms/nginx:7.2-dev": "nginx:7.2-dev",
		"docker.io/craftcms/nginx:7.1-dev": "nginx:7.1-dev",
		// "docker.io/craftcms/nginx:7.0-dev": "nginx:7.0-dev",
	}
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

			// update all of the images
			for image, name := range DockerImages {
				output.Pending("updating", name)

				// pull the image
				rdr, err := docker.ImagePull(cmd.Context(), image, types.ImagePullOptions{All: false})
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
			containers, err := docker.ContainerList(cmd.Context(), types.ContainerListOptions{All: true, Filters: filter})
			if err != nil {
				return err
			}

			// check all of the containers
			for _, c := range containers {
				// only show the site containers and proxy container
				if c.Labels[labels.Host] == "" || c.Labels[labels.Proxy] == "" {
					continue
				}

				// check the proxy container image
				if c.Labels[labels.Proxy] != "" && c.Image != fmt.Sprintf("docker.io/craftcms/nitro-proxy:%s", version.Version) {
					if err := remove(cmd.Context(), docker, c); err != nil {
						return fmt.Errorf("unable to remove container for %s: %w", strings.TrimLeft(c.Names[0], "/"), err)
					}
				}

				// if the images match, we are up to date
				if _, ok := DockerImages[c.Image]; ok {
					continue
				}

				// otherwise we need to recreate the container
				if err := remove(cmd.Context(), docker, c); err != nil {
					return fmt.Errorf("unable to remove container for %s: %w", strings.TrimLeft(c.Names[0], "/"), err)
				}

			}

			output.Info("Images updated 👍, apply changes...")

			// TODO(jasonmccallister) make this better :)
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

func remove(ctx context.Context, docker client.ContainerAPIClient, container types.Container) error {
	if err := docker.ContainerStop(ctx, container.ID, nil); err != nil {
		return err
	}

	// remove the container
	if err := docker.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
		return err
	}

	return nil
}
