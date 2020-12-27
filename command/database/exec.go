package database

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func execCreate(ctx context.Context, docker client.ContainerAPIClient, containerID string, cmds []string, show bool) (bool, error) {
	// create the exec
	e, err := docker.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd:          cmds,
	})
	if err != nil {
		return false, err
	}

	// attach to the container
	resp, err := docker.ContainerExecAttach(ctx, e.ID, types.ExecStartCheck{
		Tty: false,
	})
	defer resp.Close()

	// should we display output?
	if show {
		// show the output to stdout and stderr
		if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, resp.Reader); err != nil {
			return false, fmt.Errorf("unable to copy the output of container, %w", err)
		}
	}

	// start the exec
	if err := docker.ContainerExecStart(ctx, e.ID, types.ExecStartCheck{}); err != nil {
		return false, fmt.Errorf("unable to start the container, %w", err)
	}

	// wait for the container exec to complete
	waiting := true
	for waiting {
		resp, err := docker.ContainerExecInspect(ctx, e.ID)
		if err != nil {
			return false, err
		}

		waiting = resp.Running
	}

	return true, nil
}
