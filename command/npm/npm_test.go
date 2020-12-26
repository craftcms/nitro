// package npm

// import (
// 	"bufio"
// 	"context"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"os"
// 	"reflect"
// 	"strings"
// 	"testing"

// 	"github.com/docker/docker/api/types"
// 	"github.com/docker/docker/api/types/container"
// 	"github.com/docker/docker/api/types/filters"
// 	"github.com/docker/docker/api/types/network"
// 	"github.com/docker/docker/client"
// 	"github.com/spf13/cobra"

// 	"github.com/craftcms/nitro/labels"
// 	"github.com/craftcms/nitro/terminal"
// )

// func TestNewCommand(t *testing.T) {
// 	type args struct {
// 		docker      *mockDockerClient
// 		output      terminal.Outputer
// 		path        string
// 		environment string
// 	}
// 	tests := []struct {
// 		name                            string
// 		args                            args
// 		wantErr                         bool
// 		expectedErr                     error
// 		expectedContainerCreateRequests []types.ContainerCreateConfig
// 	}{
// 		{
// 			name: "paths with package.json do not return an error",
// 			args: args{
// 				docker: newMockDockerClient([]types.Container{
// 					{
// 						ID: "some-id",
// 					},
// 				}, []types.ImageSummary{
// 					{
// 						ID: "some-image",
// 					},
// 				}),
// 				path:        "testdata/has-package",
// 				environment: "testing-npm",
// 				output:      &spyOutputer{},
// 			},
// 			wantErr:     false,
// 			expectedErr: nil,
// 			expectedContainerCreateRequests: []types.ContainerCreateConfig{
// 				{
// 					Name: "some-container",
// 					Config: &container.Config{
// 						Image: "docker.io/library/node:14-alpine",
// 						Cmd:   []string{"npm"},
// 						Tty:   false,
// 						Labels: map[string]string{
// 							labels.Environment:        "testing-npm",
// 							labels.Type:               "npm",
// 							"com.craftcms.nitro.path": "testdata/has-package",
// 						},
// 						WorkingDir: "/home/node/app",
// 					},
// 					HostConfig: &container.HostConfig{},
// 				},
// 			},
// 		},
// 		{
// 			name: "paths without package.json return an error",
// 			args: args{
// 				docker:      newMockDockerClient(nil, nil),
// 				path:        "",
// 				environment: "testing-npm",
// 				output:      &spyOutputer{},
// 			},
// 			wantErr:     true,
// 			expectedErr: ErrNoPackageFile,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Arrange
// 			parent := &cobra.Command{}
// 			cmd := NewCommand(tt.args.docker, tt.args.output)
// 			cmd.Flags().String("environment", tt.args.environment, "testing flag")
// 			cmd.Flag("path").Value.Set(tt.args.path)
// 			cmd.Flag("keep").Value.Set("true")
// 			parent.AddCommand(cmd)

// 			// Act
// 			err := cmd.RunE(cmd, os.Args)

// 			// Assert
// 			if err != nil && tt.wantErr {
// 				if errors.Is(err, tt.expectedErr) == false {
// 					t.Errorf("expected the error to be %v, got %v instead", tt.expectedErr, err)
// 				}
// 			} else {
// 				t.Errorf("expected the error to not be nil, got: \n%v", err)
// 			}

// 			if tt.args.docker.ContainerCreateRequests != nil {
// 				if !reflect.DeepEqual(tt.args.docker.ContainerCreateRequests, tt.expectedContainerCreateRequests) {
// 					t.Errorf(
// 						"expected container create request to match\ngot:\n%v\n\nwant:\n%v",
// 						tt.args.docker.ContainerCreateRequests,
// 						tt.expectedContainerCreateRequests,
// 					)
// 				}
// 			}
// 		})
// 	}
// }

// type spyOutputer struct {
// 	infos     []string
// 	succesess []string
// 	dones     []string
// }

// func (spy spyOutputer) Info(s ...string) {
// 	spy.infos = append(spy.infos, fmt.Sprintf("%s\n", strings.Join(s, " ")))
// }

// func (spy spyOutputer) Select(r io.Reader, msg string, opts []string) (int, error) {
// 	return 0, nil
// }

// func (spy spyOutputer) Warning() {

// }

// func (spy spyOutputer) Success(s ...string) {
// 	fmt.Printf("  \u2713 %s\n", strings.Join(s, " "))
// }

// func (spy spyOutputer) Pending(s ...string) {
// 	fmt.Printf("  … %s ", strings.Join(s, " "))
// }

// func (spy spyOutputer) Done() {
// 	fmt.Print("\u2713\n")
// }

// type mockDockerClient struct {
// 	client.CommonAPIClient

// 	// filters are the filters passed to list funcs
// 	filterArgs []filters.Args

// 	// container related resources for mocking calls to the client
// 	// the fields ending in *Response are designed to capture the
// 	// requests sent to the client API.
// 	containerID              string
// 	containers               []types.Container
// 	ContainerCreateRequests  []types.ContainerCreateConfig
// 	containerCreateResponse  container.ContainerCreateCreatedBody
// 	containerStartRequests   []types.ContainerStartOptions
// 	containerRestartRequests []string

// 	images []types.ImageSummary

// 	// mockError allows us to override any func to return a method, we do not
// 	// set the error by default.
// 	mockError error
// }

// func (m *mockDockerClient) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
// 	return m.images, nil
// }

// func (m *mockDockerClient) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
// 	return m.containers, nil
// }

// func (m *mockDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error) {
// 	m.ContainerCreateRequests = append(m.ContainerCreateRequests, types.ContainerCreateConfig{
// 		Name:       containerName,
// 		Config:     config,
// 		HostConfig: hostConfig,
// 	})

// 	return m.containerCreateResponse, m.mockError
// }

// func (m *mockDockerClient) ContainerAttach(ctx context.Context, container string, options types.ContainerAttachOptions) (types.HijackedResponse, error) {
// 	rdr := ioutil.NopCloser(strings.NewReader("this"))
// 	return types.HijackedResponse{
// 		Conn:   nil,
// 		Reader: bufio.NewReader(rdr),
// 	}, nil
// }

// func (m *mockDockerClient) ContainerStart(ctx context.Context, container string, options types.ContainerStartOptions) error {
// 	return nil
// }

// // inspired by the following from the Docker docker package: https://github.com/moby/moby/blob/master/client/network_create_test.go
// func newMockDockerClient(containers []types.Container, images []types.ImageSummary) *mockDockerClient {
// 	return &mockDockerClient{
// 		images: images,
// 	}
// }
