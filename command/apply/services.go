package apply

import (
	"bytes"
	"context"
	"fmt"

	"github.com/craftcms/nitro/labels"
	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func mailhog(ctx context.Context, docker client.CommonAPIClient, output terminal.Outputer, filter filters.Args, enabled bool, networkID, env string) error {
	// add the filter for mailhog
	filter.Add("label", labels.Type+"=mailhog")
	defer filter.Add("label", labels.Type+"=mailhog")

	switch enabled {
	case true:
		// get a list of containers
		containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
		if err != nil {
			return err
		}

		if len(containers) == 0 {
			output.Pending("creating mailhog service")
		}

		// if there is no container, create it
		switch len(containers) {
		case 0:
			// pull the mailhog image
			rdr, err := docker.ImagePull(ctx, "docker.io/mailhog/mailhog", types.ImagePullOptions{})
			if err != nil {
				return err
			}

			buf := &bytes.Buffer{}
			if _, err := buf.ReadFrom(rdr); err != nil {
				return fmt.Errorf("unable to read the output from pulling the image, %w", err)
			}

			// configure the service
			smtpPort, err := nat.NewPort("tcp/udp", "1025")
			if err != nil {
				output.Warning()
				return fmt.Errorf("unable to create the port, %w", err)
			}
			httpPort, err := nat.NewPort("tcp", "8025")

			if err != nil {
				output.Warning()
				return fmt.Errorf("unable to create the port, %w", err)
			}

			containerConfig := &container.Config{
				Image: "docker.io/mailhog/mailhog",
				Labels: map[string]string{
					labels.Environment: env,
					labels.Type:        "mailhog",
				},
				ExposedPorts: nat.PortSet{
					smtpPort: struct{}{},
					httpPort: struct{}{},
				},
			}

			hostconfig := &container.HostConfig{
				PortBindings: map[nat.Port][]nat.PortBinding{
					smtpPort: {
						{
							HostIP:   "127.0.0.1",
							HostPort: "1025",
						},
					},
					httpPort: {
						{
							HostIP:   "127.0.0.1",
							HostPort: "8025",
						},
					},
				},
			}

			networkConfig := &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					env: {
						NetworkID: networkID,
					},
				},
			}

			// create the container
			if _, err := createContainer(ctx, docker, containerConfig, hostconfig, networkConfig, "mailhog"); err != nil {
				output.Warning()
				return err
			}

			output.Done()
		default:
			// check if the container is running
			if containers[0].State != "running" {
				output.Pending("starting mailhog")

				// start the container
				if err := docker.ContainerStart(ctx, containers[0].ID, types.ContainerStartOptions{}); err != nil {
					output.Warning()
					break
				}

				output.Done()

				break
			}

			output.Success("mailhog ready")
		}
	default:
		// check if there is an existing container for mailhog
		containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
		if err != nil {
			return err
		}

		// if there are no containers, we can stop here
		if len(containers) == 0 {
			return nil
		}

		// if we have a container, we need to remove it
		output.Pending("removing mailhog")

		// set the container id
		id := containers[0].ID

		// stop the container
		if err := docker.ContainerStop(ctx, id, nil); err != nil {
			output.Warning()
			output.Info(err.Error())
		}

		// remove the container
		if err := docker.ContainerRemove(ctx, id, types.ContainerRemoveOptions{RemoveVolumes: true}); err != nil {
			output.Warning()
			output.Info(err.Error())
		}
	}

	return nil
}
