package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/client"
	"github.com/fatih/color"
)

// ErrDockerPing is the error when we are aunable to reach to docker client
var ErrDockerPing = fmt.Errorf("docker does not appear to be running")

// Client represents a Nitro CLI
type Client struct {
	docker client.CommonAPIClient

	// color output
	infoOut *color.Color
	errOut  *color.Color
}

func (cli Client) Error(s ...string) {
	msg := strings.Join(s, " ")

	switch cli.errOut {
	case nil:
		fmt.Printf("%s\n", msg)
	default:
		cli.errOut.Printf("%s\n", msg)
	}
}

func (cli Client) SubError(s ...string) {
	msg := strings.Join(s, " ")

	switch cli.errOut {
	case nil:
		fmt.Printf("  ==> %s\n", msg)
	default:
		cli.errOut.Printf("  ==> %s\n", msg)
	}
}

func (cli Client) Info(s ...string) {
	msg := strings.Join(s, " ")

	switch cli.infoOut {
	case nil:
		fmt.Printf("%s\n", msg)
	default:
		cli.infoOut.Printf("%s\n", msg)
	}
}

func (cli Client) SubInfo(s ...string) {
	msg := strings.Join(s, " ")

	switch cli.infoOut {
	case nil:
		fmt.Printf("  ==> %s\n", msg)
	default:
		cli.infoOut.Printf("  ==> %s\n", msg)
	}
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
		docker:  docker,
		infoOut: color.New(color.FgCyan),
		errOut:  color.New(color.FgRed),
	}

	return cli, nil
}
