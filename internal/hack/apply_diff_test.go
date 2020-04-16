package hack

import (
	"reflect"
	"testing"

	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

func TestMountDiffActions(t *testing.T) {
	type args struct {
		name     string
		attached []config.Mount
		file     []config.Mount
	}
	tests := []struct {
		name    string
		args    args
		want    []nitro.Action
		wantErr bool
	}{
		{
			name: "additional mounts returns add actions",
			args: args{
				name: "somemachine",
				attached: []config.Mount{
					{
						Source: "./testdata/mounts/attached",
						Dest:   "already",
					},
				},
				file: []config.Mount{
					{
						Source: "./testdata/mounts/attached",
						Dest:   "/already/attached",
					},
					{
						Source: "./testdata/mounts",
						Dest:   "/some/destination",
					},
				},
			},
			want: []nitro.Action{
				{
					Type:       "mount",
					UseSyscall: false,
					Args:       []string{"mount", "./testdata/mounts", "somemachine:/some/destination"},
				},
			},
			wantErr: false,
		},
		{
			name: "removed mounts returns remove actions",
			args: args{
				name: "somemachine",
				attached: []config.Mount{
					{
						Source: "./testdata/mounts/attached",
						Dest:   "already",
					},
					{
						Source: "./testdata/mounts",
						Dest:   "/some/destination",
					},
				},
				file: []config.Mount{
					{
						Source: "./testdata/mounts/attached",
						Dest:   "/already/attached",
					},
				},
			},
			want: []nitro.Action{
				{
					Type:       "umount",
					UseSyscall: false,
					Args:       []string{"umount", "somemachine:/some/destination"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MountDiffActions(tt.args.name, tt.args.attached, tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("MountDiffActions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MountDiffActions() got = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}
