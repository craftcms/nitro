package apply

import (
	"bytes"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/terminal"
)

var (
	// DynamoDBImage is the image to use for the dynamodb
	DynamoDBImage = "amazon/dynamodb-local:latest"
)

func mailhog(ctx context.Context, docker client.CommonAPIClient, output terminal.Outputer, enabled bool, networkID string) (string, error) {
	// add the filter for mailhog
	filter := filters.NewArgs()
	filter.Add("label", labels.Type+"=mailhog")

	switch enabled {
	case true:
		// get a list of containers
		containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
		if err != nil {
			return "", err
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
				return "", err
			}

			buf := &bytes.Buffer{}
			if _, err := buf.ReadFrom(rdr); err != nil {
				return "", fmt.Errorf("unable to read the output from pulling the image, %w", err)
			}

			// configure the service
			smtpPort, err := nat.NewPort("tcp/udp", "1025")
			if err != nil {
				output.Warning()
				return "", fmt.Errorf("unable to create the port, %w", err)
			}
			httpPort, err := nat.NewPort("tcp", "8025")

			if err != nil {
				output.Warning()
				return "", fmt.Errorf("unable to create the port, %w", err)
			}

			containerConfig := &container.Config{
				Image: "docker.io/mailhog/mailhog",
				Labels: map[string]string{
					labels.Type: "mailhog",
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
					"nitro-network": {
						NetworkID: networkID,
					},
				},
			}

			// create the container
			resp, err := docker.ContainerCreate(ctx, containerConfig, hostconfig, networkConfig, nil, "mailhog.service.nitro")
			if err != nil {
				return "", fmt.Errorf("unable to create the container, %w", err)
			}

			// start the container
			if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
				return "", fmt.Errorf("unable to start the container, %w", err)
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
			return "", err
		}

		// if there are no containers, we can stop here
		if len(containers) == 0 {
			return "", nil
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

	return "", nil
}

func dynamodb(ctx context.Context, docker client.CommonAPIClient, output terminal.Outputer, enabled bool, networkID string) (string, error) {
	// add the filter for dynamodb
	filter := filters.NewArgs()
	filter.Add("label", labels.Type+"=dynamodb")

	// is the service enabled
	if enabled {
		// get a list of containers
		containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
		if err != nil {
			return "", err
		}

		// if there is not a container create it
		if len(containers) == 0 {
			// pull the mailhog image
			rdr, err := docker.ImagePull(ctx, DynamoDBImage, types.ImagePullOptions{})
			if err != nil {
				return "", err
			}

			buf := &bytes.Buffer{}
			if _, err := buf.ReadFrom(rdr); err != nil {
				return "", fmt.Errorf("unable to read the output from pulling the image, %w", err)
			}

			port, err := nat.NewPort("tcp", "8000")
			if err != nil {
				return "", err
			}

			// create the container
			resp, err := docker.ContainerCreate(ctx, &container.Config{
				Image: DynamoDBImage,
				Labels: map[string]string{
					labels.Type: "dynamodb",
				},
				ExposedPorts: nat.PortSet{
					port: struct{}{},
				},
				Cmd: []string{"-jar", "DynamoDBLocal.jar", "-sharedDb", "-dbPath", "."},
			}, &container.HostConfig{
				PortBindings: map[nat.Port][]nat.PortBinding{
					port: {
						{
							HostIP:   "127.0.0.1",
							HostPort: "8000",
						},
					},
				},
			}, &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					"nitro-network": {
						NetworkID: networkID,
					},
				},
			}, nil, "dynamodb.service.nitro")
			if err != nil {
				return "", fmt.Errorf("unable to create the container, %w", err)
			}

			// start the container
			if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
				return "", fmt.Errorf("unable to start the container, %w", err)
			}
		}
	}

	return "", nil
}
