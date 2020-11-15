package client

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
)

// Client represents a Nitro CLI
type Client struct {
	docker client.CommonAPIClient
	//out    output.Outputer
}

// NewClient creates a default docker client using the current environment.
func NewClient() (*Client, error) {
	docker, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	// check if we can talk to the docker api
	if _, err := docker.Ping(context.TODO()); err != nil {
		return nil, fmt.Errorf("docker does not appear to be running")
	}

	cli := &Client{
		docker: docker,
		//out:    output.New(),
	}

	return cli, nil
}
