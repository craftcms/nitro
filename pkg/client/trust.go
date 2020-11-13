package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

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

	// make sure there are at least some containers
	if len(containers) == 0 {
		return fmt.Errorf("unable to find the container for the proxy")
	}

	proxy := containers[0]
	stream, err := cli.Exec(ctx, proxy.ID, []string{"cat", "/data/caddy/pki/authorities/local/root.crt"})
	if err != nil {
		return fmt.Errorf("unableto exec the trust command, %w", err)
	}
	defer stream.Close()

	cert, err := ioutil.ReadAll(stream.Reader)
	if err != nil {
		return fmt.Errorf("unable to read output from exec, %w", err)
	}

	fmt.Println(strings.TrimSpace(string(cert)))

	return nil
}
