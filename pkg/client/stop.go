package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/craftcms/nitro/pkg/output"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// Stop is used to stop all containers related to a specific environment
func (cli *Client) Stop(ctx context.Context, name string, args []string) error {
	// get all the containers using a filter, we only want to stop containers which
	// have the label com.craftcms.nitro.machine=name
	filter := filters.NewArgs()
	filter.Add("label", "com.craftcms.nitro.environment="+name)
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
	if err != nil {
		return fmt.Errorf("unable to get a list of the containers, %w", err)
	}

	// if there are no containers, were done
	if len(containers) == 0 {
		output.Error("No containers are running in the", name, "environment")

		return nil
	}

	output.Info("Stopping down environment for", name)

	// stop each environment container
	for _, c := range containers {
		n := strings.TrimLeft(c.Names[0], "/")
		output.SubInfo("stopping container for", n)

		if err := cli.docker.ContainerStop(ctx, c.ID, nil); err != nil {
			errWrap := fmt.Errorf("unable to stop container %s: %w", n, err)

			output.SubError(errWrap.Error())

			return errWrap
		}
	}

	output.Info("Development environment for", name, "shutdown")

	return nil
}
