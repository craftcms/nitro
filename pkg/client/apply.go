package client

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/craftcms/nitro/internal/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/mitchellh/go-homedir"
)

// Apply is used to create a
func (cli *Client) Apply(ctx context.Context, env string, cfg config.Config) error {
	// get the network for the environment
	var networkID string

	// create a filter for the network
	filter := filters.NewArgs()
	filter.Add("label", "com.craftcms.nitro.environment="+env)

	cli.out.Info(fmt.Sprintf("Looking for the %s network", env))

	// find networks
	networks, err := cli.docker.NetworkList(ctx, types.NetworkListOptions{Filters: filter})
	if err != nil {
		return fmt.Errorf("unable to list the docker networks\n%w", err)
	}
	for _, n := range networks {
		if n.Name == env {
			networkID = n.ID
		}
	}

	// if the network is not found
	if networkID == "" {
		return fmt.Errorf("unable to find the network for %s", env)
	}

	cli.out.Info("  ==> using network", networkID)

	// get the users home dir
	home, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("unable to get the users home directory, %w", err)
	}

	// TODO(jasonmccallister) get all of the sites, their local path, the php version, and the type of project (nginx or PHP-FPM)
	cli.out.Info("Checking for existing sites")
	for _, site := range cfg.Sites {
		// add the site filter
		filter.Add("label", "com.craftcms.nitro.site="+site.Hostname)

		// TODO(jasonmccallister) make the php version dynamic based on the site
		image := fmt.Sprintf("docker.io/craftcms/php-fpm:%s-dev", "7.4")

		containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{
			All:     true,
			Filters: filter,
		})
		if err != nil {
			return fmt.Errorf("error getting a list of containers")
		}

		var containerID string
		switch len(containers) {
		case 0:
			// TODO(jasonmccallister) make this dynamic
			sourcePath := "~/dev/plugins-dev"
			if site.Hostname == "extendingcaddy.nitro" {
				sourcePath = "~/dev/extendingcaddy"
			}

			// TODO get the complete file path
			if strings.Contains(sourcePath, "~") {
				sourcePath = strings.Replace(sourcePath, "~", home, 1)
			}

			absPath, err := filepath.Abs(sourcePath)
			if err != nil {
				return fmt.Errorf("unable to get the absolute path to the site, %w", err)
			}

			if _, err := cli.docker.ImagePull(ctx, image, types.ImagePullOptions{All: false}); err != nil {
				return fmt.Errorf("unable to pull the image, %w", err)
			}

			// create the container
			resp, err := cli.docker.ContainerCreate(
				ctx,
				&container.Config{
					Image: image,
					Labels: map[string]string{
						"com.craftcms.nitro.environment": env,
						"com.craftcms.nitro.host":        site.Hostname,
					},
				},
				&container.HostConfig{
					Mounts: []mount.Mount{{
						Type: "bind",
						// TODO (jasonmccallister) get the source from the site
						Source: absPath,
						//Source: site.Webroot,
						Target: "/app",
					},
					},
				},
				&network.NetworkingConfig{
					EndpointsConfig: map[string]*network.EndpointSettings{
						env: {
							NetworkID: networkID,
						},
					},
				},
				site.Hostname,
			)
			if err != nil {
				return fmt.Errorf("unable to create the container, %w", err)
			}

			containerID = resp.ID

			cli.out.Info(fmt.Sprintf("  ==> created container for %s", site.Hostname))
		default:
			return fmt.Errorf("container already exists")
		}

		// start the container
		if err := cli.docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
			return fmt.Errorf("unable to start the container, %w", err)
		}

		// remove the site filter
		filter.Del("label", "com.craftcms.nitro.site="+site.Hostname)
	}

	//TODO(jasonmccallister) get all of the databases, engine, version, and ports and create a container for each

	// TODO(jasonmccallister) convert the sites into a Caddy json config and send to the API

	return fmt.Errorf("not yet completed")
}
