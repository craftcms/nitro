package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/client"
)

// ErrDockerPing is the error when we are unable to reach to docker client
var ErrDockerPing = fmt.Errorf("docker does not appear to be running")

// Client represents a Nitro CLI
type Client struct {
	docker client.CommonAPIClient
}

func (cli Client) Error(s ...string) {
	fmt.Printf("%s\n", strings.Join(s, " "))
}

func (cli Client) Info(s ...string) {
	fmt.Printf("%s\n", strings.Join(s, " "))
}

// func (cli Client) SubInfo(s ...string) {
// 	fmt.Printf("  \u2713 %s\n", strings.Join(s, " "))
// }

func (cli Client) InfoSuccess(s ...string) {
	fmt.Printf("  \u2713 %s\n", strings.Join(s, " "))
}

func (cli Client) InfoPending(s ...string) {
	fmt.Printf("  … %s ", strings.Join(s, " "))
}

func (cli Client) InfoDone() {
	fmt.Print("\u2713\n")
}

// NewClient creates a default docker client using the current environment.
func NewClient() (*Client, error) {
	docker, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	// check if we can talk to the docker api
	if _, err := docker.Ping(context.TODO()); err != nil {
		return nil, ErrDockerPing
	}

	cli := &Client{
		docker: docker,
		//infoOut: color.New(color.FgCyan),
		//errOut:  color.New(color.FgRed),
	}

	return cli, nil
}
