package match

import (
	"testing"

	"github.com/docker/docker/api/types"
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
