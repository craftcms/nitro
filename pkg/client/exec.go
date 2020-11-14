package client

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/docker/docker/api/types"
)

// Exec is used to execute a command in a container and retreive the response. If there is an issue creating
// the exec for the container it will return an error. The func caller is responsible for closing the reader.
func (cli *Client) Exec(ctx context.Context, containerID string, cmd []string) ([]byte, error) {
	// create an exec for the container
	exec, err := cli.docker.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		Cmd:          cmd,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create an execution for container, %w", err)
	}

	// attach to the container
	stream, err := cli.docker.ContainerExecAttach(ctx, exec.ID, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		Cmd:          cmd,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to attach to container, %w", err)
	}
	defer stream.Close()

	// read the stream content
	bytes, err := ioutil.ReadAll(stream.Reader)
	if err != nil || len(bytes) == 0 {
		return nil, fmt.Errorf("unable to read the content from container, %w", err)
	}

	return bytes, nil
}
