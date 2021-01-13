package apply

import (
	"bytes"
	"context"
	"fmt"

	"github.com/craftcms/nitro/pkg/labels"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var (
	// MinioImage is the image to use for the minio container
	MinioImage = "docker.io/minio/minio:latest"
	// MinioHost is the name for the container
	MinioHost = "minio.service.nitro"
)

type Service struct {
	Hostname string
	Port     string
}

func minio(ctx context.Context, docker client.CommonAPIClient, enabled bool, networkID string) (string, string, error) {
	// add the filter
	filter := filters.NewArgs()
	filter.Add("label", labels.Type+"=minio")

	if enabled {
		// get a list of containers
		containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
		if err != nil {
			return "", "", err
		}

		if len(containers) == 0 {
			// pull the image
			rdr, err := docker.ImagePull(ctx, MinioImage, types.ImagePullOptions{})
			if err != nil {
				return "", "", err
			}

			buf := &bytes.Buffer{}
			if _, err := buf.ReadFrom(rdr); err != nil {
				return "", "", fmt.Errorf("unable to read the output from pulling the image, %w", err)
			}

			// check if the volume needs to be created
			volumes, err := docker.VolumeList(ctx, filter)
			if err != nil {
				return "", "", fmt.Errorf("unable to list volumes, %w", err)
			}

			var volume types.Volume
			// if there are no volumes, create one
			switch len(volumes.Volumes) {
			case 1:
				volume = *volumes.Volumes[0]
			default:
				resp, err := docker.VolumeCreate(ctx, volumetypes.VolumeCreateBody{
					Driver: "local",
					Labels: map[string]string{
						labels.Nitro: "true",
						labels.Type:  "minio",
					},
				})
				if err != nil {
					return "", "", err
				}

				volume = resp
			}

			// configure the service
			port, err := nat.NewPort("tcp", "9000")
			if err != nil {
				return "", "", fmt.Errorf("unable to create the port, %w", err)
			}

			containerConfig := &container.Config{
				Image: MinioImage,
				Labels: map[string]string{
					labels.Nitro: "true",
					labels.Type:  "minio",
				},
				ExposedPorts: nat.PortSet{
					port: struct{}{},
				},
				Env: []string{"MINIO_ROOT_USER=nitro", "MINIO_ROOT_PASSWORD=nitropassword"},
				Cmd: []string{"minio", "server", "/data"},
			}

			hostconfig := &container.HostConfig{
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeVolume,
						Source: volume.Name,
						Target: "/data",
					},
				},
				PortBindings: map[nat.Port][]nat.PortBinding{
					port: {
						{
							HostIP:   "127.0.0.1",
							HostPort: "9000",
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
			resp, err := docker.ContainerCreate(ctx, containerConfig, hostconfig, networkConfig, nil, MinioHost)
			if err != nil {
				return "", "", fmt.Errorf("unable to create the container, %w", err)
			}

			// start the container
			if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
				return "", "", fmt.Errorf("unable to start the container, %w", err)
			}

			return resp.ID, MinioHost, nil
		}

		// start the container
		if err := docker.ContainerStart(ctx, containers[0].ID, types.ContainerStartOptions{}); err != nil {
			return "", "", fmt.Errorf("unable to start the container, %w", err)
		}

		return containers[0].ID, MinioHost, nil
	}

	return "", "", nil
}
