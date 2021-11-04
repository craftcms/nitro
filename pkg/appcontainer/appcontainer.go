package appcontainer

import (
	"context"
	"fmt"
	"os"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

var (
	// Image is the image used for apps with the PHP version
	Image = "docker.io/craftcms/nitro:%s"
)

func Exists(ctx context.Context, docker client.CommonAPIClient, app config.App, cfg *config.Config) (bool, error) {
	// create the filter
	filter := filters.NewArgs()
	filter.Add("label", containerlabels.Host+"="+app.GetHostname())

	// look for a container for the matching app
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return false, fmt.Errorf("error getting a list of containers")
	}

	if len(containers) > 0 {
		return true, nil
	}

	return false, nil
}

func Start(ctx context.Context, docker client.CommonAPIClient, containerId string) (string, error) {

}

func StartOrCreate(
	ctx context.Context,
	docker client.CommonAPIClient,
	homeDir string,
	networkID string,
	app config.App,
	cfg *config.Config,
) (string, error) {
	// check if nitro development is defined and override the image
	_, ok := os.LookupEnv("NITRO_DEVELOPMENT")
	if ok && app.Dockerfile == false {
		Image = "craftcms/nitro:%s"
	}

	// create the filter
	filter := filters.NewArgs()
	filter.Add("label", containerlabels.Host+"="+app.GetHostname())

	// look for a container for the matching app
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return "", fmt.Errorf("error getting a list of containers")
	}

	// if there are no containers we need to create one
	if len(containers) == 0 {
		return "", fmt.Errorf("not yet implemented")
	}

	return "", fmt.Errorf("not yet implemented")
}
