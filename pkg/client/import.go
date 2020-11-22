package client

import (
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// Import is used to copy a local file into a container at a given path. It automatically enables
// overwriting directories with files. This is used for the `db import` commands.
func (cli *Client) Import(ctx context.Context, containerID string, path string, rdr io.Reader) error {
	cli.InfoPending("preparing import")

	// copy the file to the container
	if err := cli.docker.CopyToContainer(ctx, containerID, path, rdr, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true}); err != nil {
		return fmt.Errorf("unable to copy file into the container, %w", err)
	}

	// TODO(jasonmccallister) run the import command based on the type of database

	cli.InfoDone()

	return nil
}
