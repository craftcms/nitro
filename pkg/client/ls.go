package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// LS is used to return a list of containers related to a specific environment
func (cli *Client) LS(ctx context.Context, name string, args []string) (map[string]string, error) {
	opts := make(map[string]string)

	// get all the containers using a filter, we only want to start containers which
	// have the label com.craftcms.nitro.environment=name
	filter := filters.NewArgs()
	filter.Add("label", EnvironmentLabel+"="+name)
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
	if err != nil {
		return opts, fmt.Errorf("unable to get a list of the containers, %w", err)
	}

	// if there are no containers, were done
	if len(containers) == 0 {
		return opts, ErrNoContainers
	}

	// get each of the containers
	for _, c := range containers {
		opts[strings.TrimLeft(c.Names[0], "/")] = c.ID
	}

	return opts, nil
}
