package proxycontainer

import (
	"context"
	"fmt"

	"github.com/craftcms/nitro/pkg/labels"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

var (
	// ErrNoProxyContainer is returned when the proxy container is not found
	ErrNoProxyContainer = fmt.Errorf("unable to locate the proxy container")
)

// FindAndStart will look for the proxy container and verify the container is started. It will return the
// ErrNoProxyContainer error if it is unable to locate the proxy container. It is NOT responsible for
// creating the proxy container as that is handled in the initialize package.
func FindAndStart(ctx context.Context, docker client.ContainerAPIClient) (types.Container, error) {
	// create the filters for the proxy
	f := filters.NewArgs()
	f.Add("label", labels.Type+"=proxy")

	// check if there is an existing container for the nitro-proxy
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: f, All: true})
	if err != nil {
		return types.Container{}, fmt.Errorf("unable to list the containers: %w", err)
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if n == "nitro-proxy" || n == "/nitro-proxy" {
				// check if it is running
				if c.State != "running" {
					if err := docker.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
						return types.Container{}, fmt.Errorf("unable to start the proxy container: %w", err)
					}
				}

				// return the container
				return c, nil
			}
		}
	}

	return types.Container{}, ErrNoProxyContainer
}
