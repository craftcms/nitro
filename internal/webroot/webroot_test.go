package webroot

import (
	"testing"

	"github.com/craftcms/nitro/config"
)

func TestFindWebRoot(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "can find the web dir",
			args: args{
				path: "./testdata/good-example",
			},
			want:    "web",
			wantErr: false,
		},
		{
			name: "can find the public dir",
			args: args{
				path: "./testdata/public-example",
			},
			want:    "public",
			wantErr: false,
		},
		{
			name: "can find the public_html dir",
			args: args{
				path: "./testdata/public_html-example",
			},
			want:    "public_html",
			wantErr: false,
		},
		{
			name: "can find the www dir",
			args: args{
				path: "./testdata/www-example",
			},
			want:    "www",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Find(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForExistingMount(t *testing.T) {
	type args struct {
		mount      config.Mount
		absPath    string
		webrootDir string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "use case",
			args: args{
				mount: config.Mount{
					Source: "/Users/jasonmccallister/dev",
					Dest:   "/home/ubuntu/sites",
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
				absPath:    "/Users/someuser/dev-folder/something/nested",
				webrootDir: "web",
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
				absPath:    "/Users/someuser/dev-folder/something",
				webrootDir: "web",
			},
			want: "/home/ubuntu/dev-folder/something/web",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ForExistingMount(tt.args.mount, tt.args.absPath, tt.args.webrootDir); got != tt.want {
				t.Errorf("webrootForExistingMount() = %v, want %v", got, tt.want)
			}
		})
	}
}
