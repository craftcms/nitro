package client

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

func (cli *Client) Exec(ctx context.Context, containerID string, cmd []string) (*types.HijackedResponse, error) {
	exec, err := cli.docker.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		AttachStderr: true,
		AttachStdin:  true,
		AttachStdout: true,
		Cmd:          cmd,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create an execution for container, %w", err)
	}

	stream, err := cli.docker.ContainerExecAttach(ctx, exec.ID, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		Cmd:          cmd,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to attach to container, %w", err)
	}

	// if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, stream.Reader); err != nil {
	// 	return nil, fmt.Errorf("unable to copy the output of the container logs, %w", err)
	// }

	return &stream, nil
}
