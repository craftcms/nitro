package client

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/pkg/stdcopy"
)

// Node is used to take a direction and node version to perform actions like
// npm install or update on a directory. The version flag defaults to the LTS
// but allows users to pass in any version, that has a tag on docker hub.
func (cli *Client) Node(ctx context.Context, dir, version, action string) error {
	image := fmt.Sprintf("docker.io/library/%s:%s", "node", version)

	// pull the container
	fmt.Println("Pulling node image for version", version)
	_, err := cli.docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
	if err != nil {
		return fmt.Errorf("unable to pull the docker image, %w", err)
	}

	var cmd []string
	switch action {
	case "install":
		cmd = []string{"npm", "install"}
	default:
		cmd = []string{"npm", "update"}
	}

	fmt.Println("  ==> creating temporary container for node")

	// create the temp container
	resp, err := cli.docker.ContainerCreate(ctx,
		&container.Config{
			Image: image,
			Cmd:   cmd,
			Tty:   false,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{{
				Type:   "bind",
				Source: dir,
				Target: "/app",
			},
			}},
		nil,
		"")
	if err != nil {
		return fmt.Errorf("unable to create container\n%w", err)
	}

	fmt.Println("  ==> running node", action, "this may take a moment")
	stream, err := cli.docker.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
		Logs:   true,
	})
	if err != nil {
		return fmt.Errorf("unable to attach to container, %w", err)
	}
	defer stream.Close()

	// run the container
	if err := cli.docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("unable to start the container, %w", err)
	}

	if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, stream.Reader); err != nil {
		return fmt.Errorf("unable to copy the output of the container logs, %w", err)
	}

	fmt.Println("NPM", action, "ran successfully!")

	// remove the temp container
	fmt.Println("  ==> removing temporary container")
	if err := cli.docker.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("unable to remove the temporary container %q, %w", resp.ID, err)
	}

	fmt.Println("Cleanup completed!")

	return nil
}
