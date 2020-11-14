package client

import (
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

func (cli *Client) Import(ctx context.Context, containerID string, path string, rdr io.Reader) error {
	cli.out.Info("Copying file to container")

	if err := cli.docker.CopyToContainer(ctx, containerID, path, rdr, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true}); err != nil {
		return fmt.Errorf("unable to copy to the container, %w", err)
	}

	cli.out.Info("  ==> copy completed")

	return nil
}
