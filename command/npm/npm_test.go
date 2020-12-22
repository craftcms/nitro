package npm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/craftcms/nitro/terminal"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func TestNewCommand(t *testing.T) {
	type args struct {
		docker      client.CommonAPIClient
		output      terminal.Outputer
		path        string
		environment string
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr error
	}{
		{
			name: "paths with package.json do not return an error",
			args: args{
				docker:      newMockDockerClient(nil, nil),
				path:        "testdata/has-package",
				environment: "testing-npm",
				output:      &spyOutputer{},
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "paths without package.json return an error",
			args: args{
				docker:      nil,
				path:        "",
				environment: "testing-npm",
				output:      &spyOutputer{},
			},
			wantErr:     true,
			expectedErr: ErrNoPackageFile,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			parent := &cobra.Command{}
			cmd := NewCommand(tt.args.docker, tt.args.output)
			cmd.Flags().String("environment", tt.args.environment, "testing flag")
			cmd.Flag("path").Value.Set(tt.args.path)
			parent.AddCommand(cmd)

			// Act
			err := cmd.RunE(cmd, os.Args)

			// Assert
			if err != nil && tt.wantErr {
				if errors.Is(err, tt.expectedErr) == false {
					t.Errorf("expected the error to be %v, got %v instead", tt.expectedErr, err)
				}
			} else {
				t.Error("expected the error to not be nil")
			}
		})
	}
}

type spyOutputer struct {
	infos     []string
	succesess []string
	dones     []string
}

func (spy spyOutputer) Info(s ...string) {
	spy.infos = append(spy.infos, fmt.Sprintf("%s\n", strings.Join(s, " ")))
}

func (spy spyOutputer) Select(r io.Reader, msg string, opts []string) (int, error) {
	return 0, nil
}

func (spy spyOutputer) Warning() {

}

func (spy spyOutputer) Success(s ...string) {
	fmt.Printf("  \u2713 %s\n", strings.Join(s, " "))
}

func (spy spyOutputer) Pending(s ...string) {
	fmt.Printf("  … %s ", strings.Join(s, " "))
}

func (spy spyOutputer) Done() {
	fmt.Print("\u2713\n")
}

type mockDockerClient struct {
	client.CommonAPIClient

	// filters are the filters passed to list funcs
	filterArgs []filters.Args

	// container related resources for mocking calls to the client
	// the fields ending in *Response are designed to capture the
	// requests sent to the client API.
	containerID              string
	containers               []types.Container
	containerCreateRequests  []types.ContainerCreateConfig
	containerCreateResponse  container.ContainerCreateCreatedBody
	containerStartRequests   []types.ContainerStartOptions
	containerRestartRequests []string

	images []types.ImageSummary

	// mockError allows us to override any func to return a method, we do not
	// set the error by default.
	mockError error
}

func (m *mockDockerClient) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	return nil, fmt.Errorf("not implemented")
}

// inspired by the following from the Docker docker package: https://github.com/moby/moby/blob/master/client/network_create_test.go
func newMockDockerClient(containers []types.Container, images []types.ImageSummary) *mockDockerClient {
	return &mockDockerClient{
		images: images,
	}
}
