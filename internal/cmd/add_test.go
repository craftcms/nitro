package cmd

import (
	"testing"

	"github.com/craftcms/nitro/config"
)

func Test_generateWebroot(t *testing.T) {
	type args struct {
		mount         config.Mount
		absPath       string
		webrootDir    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "use case",
			args: args{
				mount:      config.Mount{
					Source: "/Users/jasonmccallister/dev",
					Dest: "/home/ubuntu/sites",
				},
				absPath:    "/Users/jasonmccallister/dev/someproject",
				webrootDir: "web",
			},
			want: "/home/ubuntu/sites/someproject/web",
		},
		{
			name: "returns properly subnested folders",
			args: args{
				mount: config.Mount{
					Source: "/Users/someuser/dev-folder",
					Dest:   "/home/ubuntu/dev-folder",
				},
				absPath:       "/Users/someuser/dev-folder/something/nested",
				webrootDir:    "web",
			},
			want: "/home/ubuntu/dev-folder/something/nested/web",
		},
		{
			name: "returns webroot if not nested",
			args: args{
				mount: config.Mount{
					Source: "/Users/someuser/dev-folder",
					Dest:   "/home/ubuntu/dev-folder",
				},
				absPath:       "/Users/someuser/dev-folder/something",
				webrootDir:    "web",
			},
			want: "/home/ubuntu/dev-folder/something/web",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := webrootForExistingMount(tt.args.mount, tt.args.absPath, tt.args.webrootDir); got != tt.want {
				t.Errorf("webrootForExistingMount() = %v, want %v", got, tt.want)
			}
		})
	}
}
