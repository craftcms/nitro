package imageinspect

import (
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type spyClient struct {
	ImageInspect types.ImageInspect
}

func (s spyClient) BuildCancel(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s spyClient) ImageCreate(ctx context.Context, parentReference string, options types.ImageCreateOptions) (io.ReadCloser, error) {
	//TODO implement me
	panic("implement me")
}

func (s spyClient) ImageHistory(ctx context.Context, image string) ([]image.HistoryResponseItem, error) {
	//TODO implement me
	panic("implement me")
}

func (s spyClient) ImageImport(ctx context.Context, source types.ImageImportSource, ref string, options types.ImageImportOptions) (io.ReadCloser, error) {
	//TODO implement me
	panic("implement me")
}

func (s spyClient) ImageInspectWithRaw(ctx context.Context, image string) (types.ImageInspect, []byte, error) {
	return s.ImageInspect, nil, nil
}

func (s spyClient) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	//TODO implement me
	panic("implement me")
}

func (s spyClient) ImageLoad(ctx context.Context, input io.Reader, quiet bool) (types.ImageLoadResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s spyClient) ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (io.ReadCloser, error) {
	//TODO implement me
	panic("implement me")
}

func (s spyClient) ImagePush(ctx context.Context, ref string, options types.ImagePushOptions) (io.ReadCloser, error) {
	//TODO implement me
	panic("implement me")
}

func (s spyClient) ImageRemove(ctx context.Context, image string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	//TODO implement me
	panic("implement me")
}

func (s spyClient) ImageSearch(ctx context.Context, term string, options types.ImageSearchOptions) ([]registry.SearchResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s spyClient) ImageSave(ctx context.Context, images []string) (io.ReadCloser, error) {
	//TODO implement me
	panic("implement me")
}

func (s spyClient) ImageTag(ctx context.Context, image, ref string) error {
	//TODO implement me
	panic("implement me")
}

func (s spyClient) ImagesPrune(ctx context.Context, pruneFilter filters.Args) (types.ImagesPruneReport, error) {
	//TODO implement me
	panic("implement me")
}

func (s spyClient) ImageBuild(ctx context.Context, context io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	return types.ImageBuildResponse{}, nil
}

func (s spyClient) BuildCachePrune(ctx context.Context, opts types.BuildCachePruneOptions) (*types.BuildCachePruneReport, error) {
	return &types.BuildCachePruneReport{}, nil
}

func TestInspect(t *testing.T) {
	type args struct {
		ctx    context.Context
		docker client.ImageAPIClient
		image  string
	}
	tests := []struct {
		name    string
		mock    types.ImageInspect
		args    args
		want    *Info
		wantErr bool
	}{
		{
			name: "can get the required information from an image",
			mock: types.ImageInspect{
				Config: &container.Config{
					User:       "my-user",
					WorkingDir: "/var/www/html",
					ExposedPorts: map[nat.Port]struct{}{
						"80": {},
					},
				},
			},
			args: args{
				ctx:   context.TODO(),
				image: "myapp.nitro:local",
			},
			want: &Info{
				User:             "my-user",
				WorkingDirectory: "/var/www/html",
				Ports:            []int{80},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spy := spyClient{
				ImageInspect: tt.mock,
			}

			got, err := Inspect(tt.args.ctx, spy, tt.args.image)
			if (err != nil) != tt.wantErr {
				t.Errorf("Inspect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Inspect() = %v, want %v", got, tt.want)
			}
		})
	}
}
