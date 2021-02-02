package update

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/version"
	"github.com/craftcms/nitro/pkg/config"
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
	hasUpdates bool
)

// New returns the update command for updating images on the local machine as well as the nitro-proxy container.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update nitro containers and proxy",
		Example: `  # update nitro
  nitro update`,
		PostRunE: func(cmd *cobra.Command, args []string) error {
			// if there are no updates to apply return
			if !hasUpdates {
				return nil
			}

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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// get all of the php versions from the site
			versions := make(map[string]bool)
			for _, s := range cfg.Sites {
				// check if the php version is already set otherwise set it
				if _, ok := versions[s.Version]; !ok {
					versions[s.Version] = true
				}
			}

			output.Info("Updating nitro…")

			// update all of the images
			for name, image := range dockerImages {
				// make sure this is version that is installed
				ver := versionFromName(name)
				if _, ok := versions[ver]; !ok {
					// make sure its not the proxy
					if !strings.Contains(name, "proxy") {
						continue
					}
				}

				output.Pending("downloading", name)

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
				if c.Labels[labels.Type] == "dynamodb" || c.Labels[labels.Type] == "mailhog" || c.Labels[labels.Type] == "minio" || c.Labels[labels.Type] == "redis" || c.Labels[labels.Type] == "database" {
					continue
				}

				// if its the proxy container and its up to date, don't replace it
				if c.Labels[labels.Proxy] == "true" && c.Labels[labels.ProxyVersion] == version.Version {
					continue
				}

				// if the site images match, we are up to date
				if dockerImages[imageName(c.Image)] != c.Image {
					output.Pending(strings.TrimLeft(c.Names[0], "/"), "is out of date, replacing...")

					// stop the container if it is running
					if c.State == "running" {
						if err := docker.ContainerStop(ctx, c.ID, nil); err != nil {
							output.Warning()
							return err
						}
					}

					// remove the container
					if err := docker.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{}); err != nil {
						output.Warning()
						return err
					}

					hasUpdates = true

					output.Done()
				}
			}

			if hasUpdates {
				output.Info("Images updated 👍, applying changes…")
			} else {
				output.Info("Everything is up to date 👍...")
			}

			return nil
		},
	}

	return cmd
}

// get the php version from the container image name (e.g. nginx:7.3-dev)
func versionFromName(name string) string {

	p := strings.Split(name, ":")
	v := strings.Split(p[1], "-")
	return v[0]
}

// docker.io/craftcms/nginx:7.1-dev to nginx:7.1-dev
func imageName(image string) string {
	p := strings.Split(image, "/")

	return p[len(p)-1]
}
