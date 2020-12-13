package match

import (
	"testing"

	"github.com/craftcms/nitro/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

func TestMounts(t *testing.T) {
	type args struct {
		existing []types.MountPoint
		expected map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "returns false when the existing mounts do not match the expected",
			args: args{
				existing: []types.MountPoint{
					{
						Type:        mount.TypeBind,
						Source:      "~/dev/plugins/example",
						Destination: "/app/example",
					},
					{
						Type:        mount.TypeBind,
						Source:      "~/dev/mywebsite",
						Destination: "/app",
					},
				},
				expected: map[string]string{
					"~/dev/plugins/example": "/app/example",
				},
			},
			want: false,
		},
		{
			name: "returns true when the existing mounts match the expected",
			args: args{
				existing: []types.MountPoint{
					{
						Type:        mount.TypeBind,
						Source:      "~/dev/mywebsite",
						Destination: "/app",
					},
				},
				expected: map[string]string{
					"~/dev/mywebsite": "/app",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mounts(tt.args.existing, tt.args.expected); got != tt.want {
				t.Errorf("Mounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSite(t *testing.T) {
	type args struct {
		home      string
		site      config.Site
		php       config.PHP
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
					Path: "testdata/new-site",
					PHP:  "7.4",
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
					PHP: "7.4",
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
			if got := Site(tt.args.home, tt.args.site, tt.args.php, tt.args.container); got != tt.want {
				t.Errorf("Site() = %v, want %v", got, tt.want)
			}
		})
	}
}
