package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/craftcms/nitro/pkg/sudo"
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
	fmt.Println("Retreiving CA from nitro proxy")
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

	// TODO(jasonmccallister) move this to the cmd pkg

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

	fmt.Println("  ==> saved certificate to", f.Name())

	fmt.Println("  ==> attempting to add certificate, you will be prompted for a password")
	if err := sudo.Run("security", "security", "add-trusted-cert", "-d", "-r", "trustRoot", "-k", "/Library/Keychains/System.keychain", f.Name()); err != nil {
		fmt.Println("Unable to automatically add the certificate\n")

		fmt.Println("To install the certificate, run the following command:")

		// TODO show os specific commands
		switch runtime.GOOS {
		default:
			fmt.Println(fmt.Sprintf("  sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain %s", f.Name()))
		}

		return nil
	}

	// we added it correctly
	fmt.Println("  ==> certificate added")

	// clean up
	fmt.Println("  ==> removing temporary file", f.Name())

	if err := os.Remove(f.Name()); err != nil {
		fmt.Println(" ==> unable to remove temporary file, it will be automatically removed on reboot")
	}

	fmt.Println("Certificate sucessfully added, you may need to restart your browser")

	return nil
}
