package client

import (
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

	cli := &Client{
		docker: docker,
		//out:    output.New(),
	}

	return cli, nil
}
