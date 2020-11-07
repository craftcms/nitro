package client

import "github.com/docker/docker/client"

// Client represents a Nitro CLI
type Client struct {
	docker *client.Client
}

func NewClient() (*Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &Client{docker: cli}, nil
}
