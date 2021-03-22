package composer

import (
	"context"
	"fmt"
	"os/user"
	"runtime"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// Options are used to pass container specific details to the create func
type Options struct {
	Image         string
	Commands      []string
	Labels        map[string]string
	Volume        *types.Volume
	Path          string
	NetworkConfig *network.NetworkingConfig
}

// CreateContainer will create a new container for running composer with a local path and volume for caching downloads.
func CreateContainer(ctx context.Context, docker client.CommonAPIClient, opts *Options) (container.ContainerCreateCreatedBody, error) {
	if opts == nil || opts.Image == "" || opts.Commands == nil || opts.Volume == nil || opts.Path == "" {
		return container.ContainerCreateCreatedBody{}, fmt.Errorf("invalid options provided to create the composer container")
	}

	containerUser := "www-data"
	if runtime.GOOS == "linux" {
		user, err := user.Current()
		if err != nil {
			return container.ContainerCreateCreatedBody{}, err
		}
		containerUser = fmt.Sprintf("%s:%s", user.Uid, user.Gid)
	}

	return docker.ContainerCreate(
		ctx,
		&container.Config{
			Image:      opts.Image,
			Cmd:        opts.Commands,
			Tty:        false,
			Labels:     opts.Labels,
			Entrypoint: []string{"/usr/bin/composer"},
			User:       containerUser,
		},
		&container.HostConfig{
			Binds: []string{fmt.Sprintf("%s:/app:Z", opts.Path)},
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: opts.Volume.Name,
					Target: "/tmp",
				},
			},
		},
		opts.NetworkConfig,
		nil,
		"",
	)
}
