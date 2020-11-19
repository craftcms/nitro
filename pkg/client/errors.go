package client

import "fmt"

var (
	// ErrNoContainers is used when no containers are running
	ErrNoContainers = fmt.Errorf("There are no containers running")

	// ErrNoNetwork is used when we cannot find the network
	ErrNoNetwork = fmt.Errorf("Unable to find the network")
)
