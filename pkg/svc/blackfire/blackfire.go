package blackfire

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

const (
	// Image is the image to use for the blackfire container
	Image = "docker.io/blackfire/blackfire:2"

	// Host is the hostname for the blackfire container
	Host = "blackfire.service.nitro"

	// Label is the label value used to mark a container as a "blackfire" service
	Label = "blackfire"
)

// VerifyCreated will verify that the blackfire service container exists and is started
func VerifyCreated(ctx context.Context, cli client.CommonAPIClient, networkID string, cfg config.Config, output terminal.Outputer) (string, string, error) {
	// add the filter
	filter := filters.NewArgs()
	filter.Add("label", containerlabels.Nitro+"=true")
	filter.Add("label", containerlabels.Type+"="+Label)

	// get a list of containers
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return "", "", err
	}

	// if there is not a container, create one
	if len(containers) == 0 {
		// pull the image
		r, err := cli.ImagePull(ctx, Image, types.ImagePullOptions{})
		if err != nil {
			return "", "", err
		}

		// read from the buffer to pull the image
		buf := &bytes.Buffer{}
		if _, err := buf.ReadFrom(r); err != nil {
			return "", "", fmt.Errorf("unable to read output while pulling image, %w", err)
		}

		credentials, err := cfg.GetBlackfireServerCredentials()
		if err != nil {
			return "", "", err
		}

		containerConfig := &container.Config{
			Image: Image,
			Labels: map[string]string{
				containerlabels.Nitro: "true",
				containerlabels.Type:  Label,
			},
			Env: credentials,
		}

		networkConfig := &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"nitro-network": {
					NetworkID: networkID,
				},
			},
		}

		// create the container
		resp, err := cli.ContainerCreate(ctx, containerConfig, nil, networkConfig, nil, Host)
		if err != nil {
			return "", "", fmt.Errorf("unable to create the container, %w", err)
		}

		// start the container
		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return "", "", fmt.Errorf("unable to start the container, %w", err)
		}

		return resp.ID, Host, nil
	}

	// start each of the containers, there should only be one so the final return is an error
	for _, c := range containers {
		// start the container
		if c.Status != "running" {
			if err := cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
				return "", "", fmt.Errorf("unable to start the container, %w", err)
			}
		}
	}

	return containers[0].ID, Host, nil
}

// VerifyRemoved will try verify the container is not created for the minio service. If we find any containers that are
func VerifyRemoved(ctx context.Context, cli client.CommonAPIClient, output terminal.Outputer) error {
	// add the filter
	filter := filters.NewArgs()
	filter.Add("label", containerlabels.Nitro+"=true")
	filter.Add("label", containerlabels.Type+"="+Label)

	// get a list of containers
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filter,
	})
	if err != nil {
		return err
	}

	// we are all good, nothing to do
	if len(containers) == 0 {
		return nil
	}

	timeout := time.Duration(time.Second * 30)

	// remove all of the containers
	for _, c := range containers {
		// stop the container if its running
		if c.State == "running" {
			if err := cli.ContainerStop(ctx, c.ID, &timeout); err != nil {
				return err
			}
		}

		// remove the container
		if err := cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{RemoveVolumes: true}); err != nil {
			return err
		}
	}

	return nil
}
