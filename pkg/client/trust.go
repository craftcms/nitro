package client

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// Trust is used to find the proxy container and get the certificate from the container
// and place into the host machine.
func (cli *Client) Trust(ctx context.Context, env string, args []string) error {
	// find the nitro proxy for the environment
	filter := filters.NewArgs()
	filter.Add("label", "com.craftcms.nitro.environment="+env)
	filter.Add("label", "com.craftcms.nitro.proxy="+env)

	// find the container, should only be one
	containers, err := cli.docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return fmt.Errorf("unable to get the list of containers, %w", err)
	}

	// make sure there is at least one container
	if len(containers) == 0 {
		return fmt.Errorf("unable to find the container for the proxy")
	}

	// get the contents of the certificate from the container
	stream, err := cli.Exec(ctx, containers[0].ID, []string{"less", "/data/caddy/pki/authorities/local/root.crt"})
	if err != nil {
		return fmt.Errorf("unable to retreive the certificate from the proxy, %w", err)
	}
	defer stream.Close()

	// read the stream content
	bytes, err := ioutil.ReadAll(stream.Reader)
	if err != nil || len(bytes) == 0 {
		return fmt.Errorf("unable to read the content from the proxy container, %w", err)
	}

	content := string(bytes)

	fmt.Print(content)

	return nil
}
