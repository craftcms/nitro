package client

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/pkg/stdcopy"
)

// Node is used to take a direction and node version to perform actions like
// npm install or update on a directory. The version flag defaults to the LTS
// but allows users to pass in any version, that has a tag on docker hub.
func (cli *Client) Node(ctx context.Context, dir, version, action string) error {
	image := fmt.Sprintf("docker.io/library/%s:%s", "node", version)

	filters := filters.NewArgs()
	filters.Add("reference", image)

	// look for the image
	images, err := cli.docker.ImageList(ctx, types.ImageListOptions{Filters: filters})
	if err != nil {
		return fmt.Errorf("unable to get a list of images, %w", err)
	}

	// if we don't have the image, pull it
	if len(images) == 0 {
		cli.Info("Pulling node image for version", version)

		rdr, err := cli.docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
		if err != nil {
			return fmt.Errorf("unable to pull the docker image, %w", err)
		}

		buf := &bytes.Buffer{}
		if _, err := buf.ReadFrom(rdr); err != nil {
			return fmt.Errorf("unable to read the output from pulling the image, %w", err)
		}
	}

	var cmd []string
	switch action {
	case "install":
		cmd = []string{"npm", "install"}
	default:
		cmd = []string{"npm", "update"}
	}

	cli.InfoPending("preparing npm")

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

	cli.InfoDone()

	cli.Info("Running npm", action)

	// attach to the container
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

	// copy the stream to stdout
	if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, stream.Reader); err != nil {
		return fmt.Errorf("unable to copy the output of the container logs, %w", err)
	}

	// remove the temp container
	cli.InfoPending("cleaning up")

	// remove the container
	if err := cli.docker.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("unable to remove the temporary container %q, %w", resp.ID, err)
	}

	cli.InfoDone()

	cli.Info("Node", action, "complete 🤘")

	return nil
}
