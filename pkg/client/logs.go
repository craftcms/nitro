package client

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
)

func (cli *Client) Logs(ctx context.Context, containerID string) (io.ReadCloser, error) {
	return cli.docker.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{})
}
