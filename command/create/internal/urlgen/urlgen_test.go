package urlgen

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGenerate(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			name: "can get the shorthand url",
			args: args{addr: "jasonmccallister/my-craft-setup"},
			want: &url.URL{
				Scheme: "https",
				Host:   "github.com",
				Path:   "/jasonmccallister/my-craft-setup/archive/HEAD.zip",
			},
			wantErr: false,
		},
		{
			name: "can get the download URL for a complete URL",
			args: args{addr: "https://github.com/jasonmccallister/my-craft-setup"},
			want: &url.URL{
				Scheme: "https",
				Host:   "github.com",
				Path:   "/jasonmccallister/my-craft-setup/archive/HEAD.zip",
			},
			wantErr: false,
		},
		{
			name: "the default repository url is returned when nothing is entered",
			args: args{addr: ""},
			want: &url.URL{
				Scheme: "https",
				Host:   "github.com",
				Path:   "/craftcms/craft/archive/HEAD.zip",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Generate(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
