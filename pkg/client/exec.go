package client

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

// Exec is used to execute a command in a container and retreive the response. If there is an issue creating
// the exec for the container it will return an error. The func caller is responsible for closing the reader.
func (cli *Client) Exec(ctx context.Context, containerID string, cmd []string) (types.HijackedResponse, error) {
	emptyResp := types.HijackedResponse{}
	exec, err := cli.docker.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		Cmd:          cmd,
	})
	if err != nil {
		return emptyResp, fmt.Errorf("unable to create an execution for container, %w", err)
	}

	stream, err := cli.docker.ContainerExecAttach(ctx, exec.ID, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		Cmd:          cmd,
	})
	if err != nil {
		return emptyResp, fmt.Errorf("unable to attach to container, %w", err)
	}

	return stream, nil
}
