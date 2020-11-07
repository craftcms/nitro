package client

import (
	"context"
	"fmt"
	"net"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/go-connections/nat"
)

func (cli *Client) Init(ctx context.Context, name string, args []string) error {
	fmt.Println("Running pre-checks on the development environment...")

	// create filters
	filter := filters.NewArgs()
	filter.Add("name", name)

	// check if the network needs to be created
	if err := cli.checkNetwork(ctx, name, filter); err != nil {
		return err
	}

	// check if the volume needs to be created
	if err := cli.checkVolume(ctx, name, filter); err != nil {
		return err
	}

	// pull the latest image from docker hub for the nitro-proxy
	// TODO(jasonmccallister) replace with the nitro CLI version and non-local image (e.g. craftcms/nitro-proxy:version)
	//_, err = cli.docker.ImagePull(ctx, "testing-caddy:latest", types.ImagePullOptions{})
	//if err != nil {
	//	return fmt.Errorf("unable to pull the nitro-proxy from docker hub, %w", err)
	//}

	// check if port 80, 443, and 5000 are available
	ports := []string{"80", "443", "5000"}
	if err := cli.checkPorts(ports); err != nil {
		return err
	}

	// check if there is an existing container for the nitro-proxy
	if err := cli.checkContainer(ctx, name, filter); err != nil {
		return err
	}

	return nil
}

func (cli *Client) checkNetwork(ctx context.Context, name string, filter filters.Args) error {
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

		return nil
	}

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

	return nil
}

func (cli *Client) checkVolume(ctx context.Context, name string, filter filters.Args) error {
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

		return nil
	}

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

	return nil
}

func (cli *Client) checkContainer(ctx context.Context, name string, filter filters.Args) error {
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
	if err != nil {
		return fmt.Errorf("unable to list the containers\n%w", err)
	}

	// since the filter is fuzzy, do an exact match (e.g. filtering for
	// `nitro-dev` will also return `nitro-dev-host`
	var skipContainer bool
	var containerID string
	for _, c := range containers {
		for _, n := range c.Names {
			if n == name {
				skipContainer = true
				containerID = c.ID
			}
		}
	}

	// check if the volume needs to be created, the nitro-proxy container handles 80 and 443 traffic routing
	if skipContainer {
		fmt.Println(" ==> Skipping proxy container creation for nitro-proxy")
	} else {
		fmt.Println(" ==> Creating proxy container for nitro-proxy")

		resp, err := cli.docker.ContainerCreate(ctx, &container.Config{Image: "testing-caddy:latest"},
			&container.HostConfig{
				PortBindings: map[nat.Port][]nat.PortBinding{
					"80": {
						{
							"0.0.0.0",
							"80",
						},
					},
					"443": {
						{
							"0.0.0.0",
							"443",
						},
					},
					"5000": {
						{
							"0.0.0.0",
							"5000",
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

		containerID = resp.ID
	}

	// start the proxy container
	if err := cli.docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("unable to start the nitro container, %w", err)
	}

	return nil
}

func (cli *Client) checkPorts(ports []string) error {
	for _, port := range ports {
		lis, err := net.Listen("tcp", ":"+port)
		if err != nil {
			return fmt.Errorf("nitro uses ports 80, 443, and 5000. It appears port %q, is already in use", port)
		}

		if err := lis.Close(); err != nil {
			return fmt.Errorf("unable to close the listener after checking the ports, %w", err)
		}
	}

	return nil
}
