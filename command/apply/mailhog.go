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
)

var (
	// MailhogImage is the image to use for the mailhog container
	MailhogImage = "docker.io/mailhog/mailhog:latest"
	MailhogHost  = "mailhog.service.nitro"
)

func mailhog(ctx context.Context, docker client.CommonAPIClient, enabled bool, networkID string) (string, string, error) {
	// add the filter
	filter := filters.NewArgs()
	filter.Add("label", labels.Type+"=mailhog")

	if enabled {
		// get a list of containers
		containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
		if err != nil {
			return "", "", err
		}

		if len(containers) == 0 {
			// pull the image
			rdr, err := docker.ImagePull(ctx, MailhogImage, types.ImagePullOptions{})
			if err != nil {
				return "", "", err
			}

			buf := &bytes.Buffer{}
			if _, err := buf.ReadFrom(rdr); err != nil {
				return "", "", fmt.Errorf("unable to read the output from pulling the image, %w", err)
			}

			// configure the service
			smtpPort, err := nat.NewPort("tcp/udp", "1025")
			if err != nil {
				return "", "", fmt.Errorf("unable to create the port, %w", err)
			}
			httpPort, err := nat.NewPort("tcp", "8025")
			if err != nil {
				return "", "", fmt.Errorf("unable to create the port, %w", err)
			}

			containerConfig := &container.Config{
				Image: MailhogImage,
				Labels: map[string]string{
					labels.Nitro: "true",
					labels.Type:  "mailhog",
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
			resp, err := docker.ContainerCreate(ctx, containerConfig, hostconfig, networkConfig, nil, MailhogHost)
			if err != nil {
				return "", "", fmt.Errorf("unable to create the container, %w", err)
			}

			// start the container
			if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
				return "", "", fmt.Errorf("unable to start the container, %w", err)
			}

			return resp.ID, MailhogHost, nil
		}

		// start the container
		if err := docker.ContainerStart(ctx, containers[0].ID, types.ContainerStartOptions{}); err != nil {
			return "", "", fmt.Errorf("unable to start the container, %w", err)
		}

		return containers[0].ID, MailhogHost, nil
	}

	return "", "", nil
}
