package match

import (
	"testing"

	"github.com/craftcms/nitro/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func TestSite(t *testing.T) {
	type args struct {
		home      string
		site      config.Site
		container types.ContainerJSON
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "path mismatches return false",
			args: args{
				home: "testdata/example-site",
				site: config.Site{
					Path:    "testdata/new-site",
					Version: "7.4",
				},
				container: types.ContainerJSON{
					Config: &container.Config{
						Image: "docker.io/craftcms/nginx:7.4-dev",
					},
				},
			},
			want: false,
		},
		{
			name: "mismatched images return false",
			args: args{
				site: config.Site{
					Version: "7.4",
				},
				container: types.ContainerJSON{
					Config: &container.Config{
						Image: "docker.io/craftcms/nginx:7.3-dev",
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Site(tt.args.home, tt.args.site, tt.args.container); got != tt.want {
				t.Errorf("Site() = %v, want %v", got, tt.want)
			}
		})
	}
}
