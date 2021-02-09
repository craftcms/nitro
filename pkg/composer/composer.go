package composer

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

// Options are used to pass container specific details to the create func
type Options struct {
	Image    string
	Commands []string
	Labels   map[string]string
	Volume   *types.Volume
	Path     string
}

// CreateContainer will create a new container for running composer with a local path and volume for caching downloads.
func CreateContainer(ctx context.Context, docker client.CommonAPIClient, opts *Options) (container.ContainerCreateCreatedBody, error) {
	if opts == nil || opts.Image == "" || opts.Commands == nil || opts.Volume == nil || opts.Path == "" {
		return container.ContainerCreateCreatedBody{}, fmt.Errorf("invalid options provided to create the composer container")
	}

	return docker.ContainerCreate(
		ctx,
		&container.Config{
			Image:  opts.Image,
			Cmd:    opts.Commands,
			Tty:    false,
			Labels: opts.Labels,
			Env:    []string{"COMPOSER_HOME=/root"},
		},
		&container.HostConfig{Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: opts.Volume.Name,
				Target: "/root",
			},
			{
				Type:   mount.TypeBind,
				Source: opts.Path,
				Target: "/app",
			},
		},
		},
		nil,
		nil,
		"",
	)
}
