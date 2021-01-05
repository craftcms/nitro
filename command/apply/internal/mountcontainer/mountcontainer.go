package mountcontainer

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/craftcms/nitro/command/apply/internal/match"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

var (
	MountImage = "docker.io/craftcms/php-fpm:%s-dev"
)

func FindOrCreate(ctx context.Context, docker client.CommonAPIClient, home, networkID string, mount config.Mount) (string, error) {
	// set filters for the container
	filter := filters.NewArgs()
	filter.Add("label", labels.Type+"=mount")

	// look for a container for the site
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return "", fmt.Errorf("error getting a list of containers")
	}

	// if there are no matching containers, we need to create
	if len(containers) == 0 {
		return create(ctx, docker, home, networkID, mount)
	}

	// there is a container, so inspect it and make sure it matched
	container := containers[0]

	// get the containers details that include environment variables
	details, err := docker.ContainerInspect(ctx, container.ID)
	if err != nil {
		return "", err
	}

	// check its config/envs
	if !match.Mount(home, mount, details) {
		fmt.Print("- updating... ")

		// stop container
		if err := docker.ContainerStop(ctx, container.ID, nil); err != nil {
			return "", err
		}

		// remove container
		if err := docker.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
			return "", err
		}

		return create(ctx, docker, home, networkID, mount)
	}

	// start the container
	if err := docker.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("unable to start container, %w", err)
	}

	return container.ID, nil
}

func create(ctx context.Context, docker client.CommonAPIClient, home, networkID string, mnt config.Mount) (string, error) {
	// create the container
	image := fmt.Sprintf(MountImage, mnt.Version)

	// pull the image
	rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
	if err != nil {
		return "", fmt.Errorf("unable to pull the image, %w", err)
	}

	buf := &bytes.Buffer{}
	if _, err := buf.ReadFrom(rdr); err != nil {
		return "", fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
	}

	p := mnt.Path
	if strings.Contains(p, "~") {
		p = strings.Replace(p, "~", home, -1)
	}

	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}

	p = filepath.Clean(abs)

	envs := mnt.AsEnvs("")

	// create the container
	resp, err := docker.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			Labels: map[string]string{
				labels.Nitro: "true",
				labels.Type:  "mount",
				labels.Path:  p,
			},
			Env: envs,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: p,
					Target: "/app",
				},
			},
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"nitro-network": {
					NetworkID: networkID,
				},
			},
		},
		nil,
		containerName(mnt),
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

func containerName(mount config.Mount) string {
	// remove the home directory
	n := strings.Replace(mount.Path, "~/", "", 1)

	// replace path separator with underscores
	n = strings.Replace(n, string(os.PathSeparator), "_", -1)

	return fmt.Sprintf("mount_%s", n)
}
