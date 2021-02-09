package customcontainer

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/craftcms/nitro/command/apply/internal/match"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/pathexists"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func StartOrCreate(ctx context.Context, docker client.CommonAPIClient, home, networkID string, c config.Container) (hostname string, err error) {
	// set filters for the container
	filter := filters.NewArgs()
	filter.Add("label", labels.Nitro+"=true")
	filter.Add("label", labels.NitroContainer+"="+c.Name)

	// look for a container for the site
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return "", fmt.Errorf("error getting a list of containers")
	}

	// if there are no containers we need to create one
	if len(containers) == 0 {
		return create(ctx, docker, home, networkID, c)
	}

	// there is a container, so inspect it and make sure it matched
	container := containers[0]

	// start the container if not running
	if container.State != "running" {
		if err := docker.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
			return "", err
		}
	}

	// get the containers details that include environment variables
	details, err := docker.ContainerInspect(ctx, container.ID)
	if err != nil {
		return "", err
	}

	// if the container is out of date
	if err := match.Container(home, c, details); err != nil {
		fmt.Println(err)
		fmt.Print("- updating… ")

		// stop container
		if err := docker.ContainerStop(ctx, container.ID, nil); err != nil {
			return "", err
		}

		// remove container
		if err := docker.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
			return "", err
		}

		return create(ctx, docker, home, networkID, c)
	}

	return container.ID, nil
}

func create(ctx context.Context, docker client.CommonAPIClient, home, networkID string, c config.Container) (string, error) {
	// create the container
	image := fmt.Sprintf("%s:%s", c.Image, c.Tag)

	// pull the image
	rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
	if err != nil {
		return "", fmt.Errorf("unable to pull the image, %w", err)
	}

	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(rdr); err != nil {
		return "", fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
	}

	// get the containers custom environment variables from the file
	var customEnvs []string
	if c.EnvFile != "" {
		// get the file
		envFilePath := filepath.Join(home, ".nitro", "."+c.Name)

		// make sure it exists
		if !pathexists.IsFile(envFilePath) {
			return "", fmt.Errorf("unable to find file: %q", envFilePath)
		}

		content, err := ioutil.ReadFile(envFilePath)
		if err != nil {
			return "", err
		}

		for _, l := range strings.Split(string(content), "\n") {
			if strings.Contains(l, "=") {
				customEnvs = append(customEnvs, l)
			}
		}
	}

	lbls := map[string]string{
		labels.Nitro:          "true",
		labels.Type:           "custom",
		labels.NitroContainer: c.Name,
	}

	config := &container.Config{
		Image:  image,
		Labels: lbls,
	}

	if len(customEnvs) > 0 {
		config.Env = customEnvs
	}

	// check for volumes and create if not found
	var mounts []mount.Mount
	if len(c.Volumes) > 0 {
		for _, v := range c.Volumes {
			// generate the volume name
			name := fmt.Sprintf("nitro_%s_%s", c.Name, strings.Replace(v, "/", "_", -1))

			// filter for the volume
			volFilter := filters.NewArgs()
			volFilter.Add("name", name)

			// check for an existing volume
			resp, err := docker.VolumeList(ctx, volFilter)
			if err != nil {
				return "", err
			}

			if len(resp.Volumes) == 0 {
				vol, err := docker.VolumeCreate(ctx, volume.VolumeCreateBody{Driver: "local", Name: name, Labels: lbls})
				if err != nil {
					return "", err
				}

				// append the mount
				mounts = append(mounts, mount.Mount{
					Type:   mount.TypeVolume,
					Source: vol.Name,
					Target: v,
				})
			}
		}
	}

	// create the expose nat port settings and host config port bindings
	portBindings := make(map[nat.Port][]nat.PortBinding)
	portSettings := make(map[nat.Port]struct{})
	if len(c.Ports) > 0 {
		for _, p := range c.Ports {
			parts := strings.Split(p, ":")
			// get the host port
			host := parts[0]
			// get the container port
			container := parts[1]

			port, err := nat.NewPort("tcp", container)
			if err != nil {
				return "", fmt.Errorf("unable to create the port, %w", err)
			}

			// create the port binding
			portBindings[port] = []nat.PortBinding{{
				HostIP:   "127.0.0.1",
				HostPort: host,
			}}

			// create the nat port setting
			portSettings[port] = struct{}{}
		}
	}

	config.ExposedPorts = portSettings

	// create the container
	resp, err := docker.ContainerCreate(
		ctx,
		config,
		&container.HostConfig{
			Mounts:       mounts,
			PortBindings: portBindings,
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"nitro-network": {
					NetworkID: networkID,
				},
			},
		},
		nil,
		fmt.Sprintf("%s.containers.nitro", c.Name),
	)
	if err != nil {
		return "", fmt.Errorf("unable to create the container, %w", err)
	}

	// start the container
	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("unable to start the container, %w", err)
	}

	return resp.ID, nil
}
