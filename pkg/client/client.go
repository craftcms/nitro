package client

import "github.com/docker/docker/client"

// Client represents a Nitro CLI
type Client struct {
	docker client.CommonAPIClient
}

// NewClient creates a default docker client using the current environment.
func NewClient() (*Client, error) {
	cli, err := client.NewEnvClient()
	return &Client{docker: cli}, err
}
