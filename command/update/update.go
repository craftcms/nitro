package update

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"

	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/command/version"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	DockerImages = map[string]string{
		"nginx:8.0-dev":                  "docker.io/craftcms/nginx:8.0-dev",
		"nginx:7.4-dev":                  "docker.io/craftcms/nginx:7.4-dev",
		"nginx:7.3-dev":                  "docker.io/craftcms/nginx:7.3-dev",
		"nginx:7.2-dev":                  "docker.io/craftcms/nginx:7.2-dev",
		"nginx:7.1-dev":                  "docker.io/craftcms/nginx:7.1-dev",
		"nginx:7.0-dev":                  "docker.io/craftcms/nginx:7.0-dev",
		"nitro-proxy:" + version.Version: "docker.io/craftcms/nitro-proxy:" + version.Version,
	}
	runApply bool
)

// New returns the update command for updating images on the local machine as well as the nitro-proxy container.
func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates Nitro containers and proxy.",
		Example: `  # update nitro
  nitro update`,
		PostRunE: func(cmd *cobra.Command, args []string) error {
			// if there are no updates to apply return
			if !runApply {
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
			debug, err := strconv.ParseBool(cmd.Flag("debug").Value.String())
			if err != nil {
				debug = false
			}

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
			for name, image := range DockerImages {
				// make sure this is version that is installed and not the proxy
				if _, ok := versions[versionFromName(name)]; !ok && !strings.Contains(name, "proxy") {
					continue
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
					output.Warning()

					return fmt.Errorf("unable to read the output while pulling image, %w", err)
				}

				output.Done()
			}

			// create a filter for nitro containers
			filter := filters.NewArgs()
			filter.Add("label", containerlabels.Nitro)

			// get a list of containers
			containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
			if err != nil {
				return err
			}

			// check all of the containers
			for _, container := range containers {
				// is this a database, service, composer, or node container?
				if container.Labels[containerlabels.Type] == "dynamodb" || container.Labels[containerlabels.Type] == "mailhog" || container.Labels[containerlabels.Type] == "minio" || container.Labels[containerlabels.Type] == "redis" || container.Labels[containerlabels.Type] == "database" {
					continue
				}

				// is this image up to date?
				if _, ok := DockerImages[shortImageName(container.Image)]; ok {
					continue
				}

				output.Pending(strings.TrimLeft(container.Names[0], "/"), "is out of date, replacing...")

				// if we are dubugging, don't actually remove or apply changes
				if !debug {
					// stop the container if it is running
					if container.State == "running" {
						if err := docker.ContainerStop(ctx, container.ID, nil); err != nil {
							output.Warning()
							return err
						}
					}

					// remove the container
					if err := docker.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
						output.Warning()
						return err
					}

					// we need to run the apply command
					runApply = true
				}

				output.Done()
			}

			// if there are changes show a apply changes prompt
			if runApply {
				output.Info("Images updated 👍, applying changes…")
			} else {
				output.Info("Everything is up to date 👍...")
			}

			return nil
		},
	}

	cmd.Flags().Bool("debug", false, "Show what will be updated without removing the container")

	return cmd
}

// docker.io/craftcms/nginx:7.4-dev => nginx:7.4-dev
func shortImageName(s string) string {
	parts := strings.Split(s, "/")

	return parts[len(parts)-1]
}

// get the php version from the container image name (e.g. nginx:7.3-dev)
func versionFromName(name string) string {
	p := strings.Split(name, ":")
	v := strings.Split(p[1], "-")
	return v[0]
}
