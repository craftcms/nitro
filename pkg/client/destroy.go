package client

import (
	"context"
	"fmt"
	"time"
)

// Destroy is used to completed remove all containers, networks, and volumes for an environment.
// is it a destructive action and will prompt a user for verification and perform a database
// backup before removing the resources.
func (cli *Client) Destroy(ctx context.Context, name string, args []string) error {
	// TODO(jasonmccallister) get all related containers
	// TODO(jasonmccallister) get all related volumes
	// TODO(jasonmccallister) get all related networks

	// if there are no containers, were done
	// get all the containers
	fmt.Println("Backing up databases")
	fmt.Println("  ==> stopping fake container name")
	time.Sleep(time.Second * 2)


	// stop all of the container
	fmt.Println("Stopping containers")
	fmt.Println("  ==> stopping fake container name")
	time.Sleep(time.Second * 2)

	// remove all containers
	fmt.Println("Removing containers")
	fmt.Println("  ==> removing fake container name")
	time.Sleep(time.Second * 2)

	// get all the volumes
	fmt.Println("Removing volumes")
	fmt.Println("  ==> removing volume name")
	time.Sleep(time.Second * 2)

	// get all the networks
	fmt.Println("Removing network")
	fmt.Println("  ==> removing network name")
	time.Sleep(time.Second * 2)

	return fmt.Errorf("not yet implemented")
}
