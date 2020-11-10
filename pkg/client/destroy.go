package client

import (
	"context"
	"fmt"
)

// Destroy is used to completed remove all containers, networks, and volumes for an environment.
// is it a destructive action and will prompt a user for verification and perform a database
// backup before removing the resources.
func (cli *Client) Destroy(ctx context.Context, name string, args []string) error {

	// if there are no containers, were done

	return fmt.Errorf("not yet implemented")
}
