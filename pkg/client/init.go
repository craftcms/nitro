package client

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	volumetypes "github.com/docker/docker/api/types/volume"
	"github.com/docker/go-connections/nat"
)

// Init is responsible for creating the networks, volumes, and proxy containers for the
// nitro docker setup. It will check for the existing networks, volumes, and containers
// instead of overwriting the containers. Init should only be used once to setup the
// development environment, which is why we safeguard the resources.
func (cli *Client) Init(ctx context.Context, name string, args []string) error {
	fmt.Println("Running pre-checks on", name, "development environment...")

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
	var networkID string
	for _, n := range networks {
		if n.Name == name {
			skipNetwork = true
			networkID = n.ID
		}
	}

	// create the network needs to be created
	if skipNetwork {
		fmt.Println("  ==> skipping network")
	} else {
		fmt.Println("  ==> creating network")

		resp, err := cli.docker.NetworkCreate(ctx, name, types.NetworkCreate{
			Driver:     "bridge",
			Attachable: true,
			Labels: map[string]string{
				"com.craftcms.nitro.environment": name,
				"com.craftcms.nitro.network":     name,
			},
		})
		if err != nil {
			return fmt.Errorf("unable to create the network, %w", err)
		}

		// set the newly created network
		networkID = resp.ID

		fmt.Println("  ==> network created for", name)
	}

	// check if the volume needs to be created
	volumes, err := cli.docker.VolumeList(ctx, filter)
	if err != nil {
		return fmt.Errorf("unable to list the docker volumes, %w", err)
	}

	// since the filter is fuzzy, do an exact match (e.g. filtering for
	// `nitro-dev` will also return `nitro-dev-host`
	var skipVolume bool
	var volume *types.Volume
	for _, v := range volumes.Volumes {
		if v.Name == name {
			skipVolume = true
			volume = v
		}
	}

	// check if the volume needs to be created
	if skipVolume {
		fmt.Println("  ==> skipping volume")
	} else {
		fmt.Println("  ==> creating volume")

		// create a volume with the same name of the machine
		resp, err := cli.docker.VolumeCreate(ctx, volumetypes.VolumesCreateBody{
			Driver: "local",
			Name:   name,
			Labels: map[string]string{
				"com.craftcms.nitro.environment": name,
				"com.craftcms.nitro.volume":      name,
			},
		})
		if err != nil {
			return fmt.Errorf("unable to create the network, %w", err)
		}

		volume = &resp

		fmt.Println("  ==> volume created for", name)
	}

	// pull the latest image from docker hub for the nitro-proxy
	// TODO(jasonmccallister) replace with the nitr o CLI version and non-local image (e.g. craftcms/nitro-proxy:version)
	// _, err = cli.docker.ImagePull(ctx, "nitro-proxy:develop", types.ImagePullOptions{})
	// if err != nil {
	// 	return fmt.Errorf("unable to pull the nitro-proxy from docker hub, %w", err)
	// }

	// create a filter for the nitro proxy
	pf := filters.NewArgs()
	pf.Add("label", "com.craftcms.nitro.proxy="+name)

	// check if there is an existing container for the nitro-proxy
	var containerID string
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{Filters: pf, All: true})
	if err != nil {
		return fmt.Errorf("unable to list the containers\n%w", err)
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if n == name || n == "/"+name {
				fmt.Println("  ==> skipping proxy")

				containerID = c.ID
			}
		}
	}

	// if we do not have a container id, it needs to be create
	if containerID == "" {
		fmt.Println("  ==> creating container for the proxy")
		resp, err := cli.docker.ContainerCreate(ctx,
			&container.Config{
				// TODO(jasonmccallister) make this dynamic based on the nitro CLI and image
				Image: "nitro-proxy:develop",
				ExposedPorts: nat.PortSet{
					"80/tcp":   struct{}{},
					"443/tcp":  struct{}{},
					"5000/tcp": struct{}{},
				},
				Labels: map[string]string{
					"com.craftcms.nitro.type":        "proxy",
					"com.craftcms.nitro.environment": name,
					"com.craftcms.nitro.proxy":       name,
				},
			},
			&container.HostConfig{
				NetworkMode: "default",
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeVolume,
						Source: volume.Name,
						Target: "/data",
					},
				},
				// TODO(jasonmccallister) make the ports for HTTP, HTTPS, and the gRPC API dynamic
				PortBindings: map[nat.Port][]nat.PortBinding{
					"80/tcp": {
						{
							HostIP:   "127.0.0.1",
							HostPort: "80",
						},
					},
					"443/tcp": {
						{
							HostIP:   "127.0.0.1",
							HostPort: "443",
						},
					},
					"5000/tcp": {
						{
							HostIP:   "127.0.0.1",
							HostPort: "5000",
						},
					},
				},
			},
			&network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{
					name: {
						NetworkID: networkID,
					},
				},
			},
			name,
		)
		if err != nil {
			return fmt.Errorf("unable to create the container\n%w", err)
		}

		containerID = resp.ID
	}

	// start the container for the proxy
	fmt.Println("  ==> starting proxy container")

	if err := cli.docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("unable to start the nitro container, %w", err)
	}

	fmt.Println("Development environment for", name, "started")

	return nil
}
