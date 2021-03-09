package dynamodb

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	// Image is the image to use for the dynamodb container
	Image = "docker.io/amazon/dynamodb-local:latest"

	// Host is the hostname for the dynamodb container
	Host = "dynamodb.service.nitro"

	// Label is the label value used to mark a container as a "dynamodb" service
	Label = "dynamodb"
)

// VerifyCreated will verify that the dynamodb service container exists and is started
func VerifyCreated(ctx context.Context, cli client.CommonAPIClient, networkID string, output terminal.Outputer) (string, string, error) {
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

		// set the nitro env overrides
		httpPort := "8000"
		if os.Getenv("NITRO_DYNAMODB_PORT") != "" {
			httpPort = os.Getenv("NITRO_DYNAMODB_PORT")
		}

		httpPortNat, err := nat.NewPort("tcp", "8000")
		if err != nil {
			return "", "", fmt.Errorf("unable to create the port, %w", err)
		}

		containerConfig := &container.Config{
			Image: Image,
			Labels: map[string]string{
				containerlabels.Nitro: "true",
				containerlabels.Type:  Label,
			},
			ExposedPorts: nat.PortSet{
				httpPortNat: struct{}{},
			},
			Cmd: []string{"-jar", "DynamoDBLocal.jar", "-sharedDb", "-dbPath", "."},
		}

		hostconfig := &container.HostConfig{
			PortBindings: map[nat.Port][]nat.PortBinding{
				httpPortNat: {
					{
						HostIP:   "127.0.0.1",
						HostPort: httpPort,
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
		resp, err := cli.ContainerCreate(ctx, containerConfig, hostconfig, networkConfig, nil, Host)
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
		if err := cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
		}); err != nil {
			return err
		}
	}

	return nil
}
