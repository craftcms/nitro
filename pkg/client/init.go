package client

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/go-connections/nat"
)

// Init is responsible for creating the networks, volumes, and proxy containers for the
// nitro docker setup. It will check for the existing networks, volumes, and containers
// instead of overwriting the containers. Init should only be used once to setup the
// development environment, which is why we safeguard the resources.
func (cli *Client) Init(ctx context.Context, name string, args []string) error {
	fmt.Println("Running pre-checks on the development environment...")

	// create filters for the development environment
	filter := filters.NewArgs()
	filter.Add("name", name)

	// check if the network needs to be created
	networks, err := cli.docker.NetworkList(ctx, types.NetworkListOptions{Filters: filter})
	if err != nil {
		return fmt.Errorf("unable to list the docker networks, %w", err)
	}

	// since the filter is fuzzy, do an exact match (e.g. filtering for
	// `nitro-dev` will also return `nitro-dev-host`
	var skipNetwork bool
	for _, n := range networks {
		if n.Name == name {
			skipNetwork = true
		}
	}

	// create the network needs to be created
	if skipNetwork {
		fmt.Println(" ==> Skipping network creation for", name)
	} else {
		fmt.Println(" ==> Creating network for", name)

		resp, err := cli.docker.NetworkCreate(ctx, name, types.NetworkCreate{
			Driver:     "bridge",
			Attachable: true,
			Labels: map[string]string{
				"nitro": name,
			},
		})
		if err != nil {
			return fmt.Errorf("unable to create the network, %w", err)
		}

		fmt.Println(" ==> network created with id", resp.ID)
	}

	// check if the volume needs to be created
	volumes, err := cli.docker.VolumeList(ctx, filter)
	if err != nil {
		return fmt.Errorf("unable to list the docker volumes, %w", err)
	}

	// since the filter is fuzzy, do an exact match (e.g. filtering for
	// `nitro-dev` will also return `nitro-dev-host`
	var skipVolume bool
	for _, volume := range volumes.Volumes {
		if volume.Name == name {
			skipVolume = true
		}
	}

	// check if the volume needs to be created
	if skipVolume {
		fmt.Println(" ==> Skipping volume creation for", name)
	} else {
		fmt.Println(" ==> Creating volume for", name)

		// create a volume with the same name of the machine
		resp, err := cli.docker.VolumeCreate(ctx, volumetypes.VolumesCreateBody{
			Driver: "local",
			Name:   name,
			Labels: map[string]string{
				"nitro": name,
			},
		})
		if err != nil {
			return fmt.Errorf("unable to create the network, %w", err)
		}

		fmt.Println(" ==> volume created with name", resp.Name)
	}

	// pull the latest image from docker hub for the nitro-proxy
	// TODO(jasonmccallister) replace with the nitro CLI version and non-local image (e.g. craftcms/nitro-proxy:version)
	//_, err = cli.docker.ImagePull(ctx, "testing-caddy:latest", types.ImagePullOptions{})
	//if err != nil {
	//	return fmt.Errorf("unable to pull the nitro-proxy from docker hub, %w", err)
	//}

	// check if there is an existing container for the nitro-proxy
	pf := filters.NewArgs()
	pf.Add("name", "nitro-proxy")
	if err := cli.checkContainer(ctx, "nitro-proxy", pf); err != nil {
		return err
	}

	return nil
}

func (cli *Client) checkContainer(ctx context.Context, name string, filter filters.Args) error {
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter, All: true})
	if err != nil {
		return fmt.Errorf("unable to list the containers\n%w", err)
	}

	// since the filter is fuzzy, do an exact match (e.g. filtering for
	// `nitro-dev` will also return `nitro-dev-host`
	var skipContainer bool
	var containerID string
	for _, c := range containers {
		for _, n := range c.Names {
			if n == name || n == fmt.Sprintf("/%s", name) {
				skipContainer = true
				containerID = c.ID
			}
		}
	}

	// check if the volume needs to be created, the nitro-proxy container handles traffic routing
	if skipContainer {
		fmt.Println(" ==> Skipping proxy container creation for nitro-proxy")

		return cli.startContainer(ctx, containerID)
	}

	fmt.Println(" ==> Creating proxy container for nitro-proxy")

	resp, err := cli.docker.ContainerCreate(ctx,
		&container.Config{Image: "testing-caddy:latest"},
		&container.HostConfig{
			// TODO(jasonmccallister) make the ports for HTTP, HTTPS, and the gRPC API dynamic
			PortBindings: map[nat.Port][]nat.PortBinding{
				"80": {
					{
						HostIP:   "localhost",
						HostPort: "80",
					},
				},
				"443": {
					{
						HostIP:   "localhost",
						HostPort: "443",
					},
				},
				"5000": {
					{
						HostIP:   "localhost",
						HostPort: "5000",
					},
				},
			},
		},
		&network.NetworkingConfig{},
		"nitro-proxy",
	)
	if err != nil {
		return fmt.Errorf("unable to create the nitro container\n%w", err)
	}

	// start the proxy container
	return cli.startContainer(ctx, resp.ID)
}

func (cli *Client) startContainer(ctx context.Context, containerID string) error {
	if err := cli.docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("unable to start the nitro container, %w", err)
	}

	return nil
}
