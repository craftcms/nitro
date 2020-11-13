package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

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
	cli.out.Info("Retreiving CA from nitro proxy")
	content, err := cli.Exec(ctx, containers[0].ID, []string{"less", "/data/caddy/pki/authorities/local/root.crt"})
	if err != nil {
		return fmt.Errorf("unable to retreive the certificate from the proxy, %w", err)
	}

	// remove special characters from the output
	var stop int
	for i, s := range content {
		if s != 0 && s != 1 {
			stop = i + 1
			break
		}
	}

	// create a temp file
	f, err := ioutil.TempFile(os.TempDir(), "nitro-cert")
	if err != nil {
		return fmt.Errorf("unable to create a temporary file, %w", err)
	}

	// write the certificate to the temporary file
	if _, err := f.Write(content[stop+1:]); err != nil {
		return fmt.Errorf("unable to write the certificate to the temporary file, %w", err)
	}
	defer f.Close()

	cli.out.Info("  ==> saved certificate to", f.Name())

	if runtime.GOOS == "darwin" {
		cli.out.Info("To install the certificate, run the following command:")
		cli.out.Info(fmt.Sprintf("  sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain %s", f.Name()))
	}

	return nil
}
