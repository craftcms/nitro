package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/craftcms/nitro/pkg/config"
)

func TestApp_GetHostname(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	home := filepath.Join(wd, "testdata")

	type fields struct {
		Config     string
		Dockerfile bool
		Hostname   string
		Aliases    []string
		Path       string
		Webroot    string
		Extensions []string
		Xdebug     bool
		Blackfire  bool
		Database   struct {
			Engine  string `yaml:"engine,omitempty"`
			Version string `yaml:"version,omitempty"`
		}
		HomeDir string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Can get the hostname from the app config",
			fields: fields{
				Hostname: "mysite.nitro",
				HomeDir:  home,
			},
			want: "mysite.nitro",
		},
		{
			name: "Can get the hostname from the config file path",
			fields: fields{
				Config:  "~/path-with-config/nitro.yaml",
				HomeDir: home,
			},
			want: "custom-hostname-from-file.nitro",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := config.App{
				Config:     tt.fields.Config,
				Dockerfile: tt.fields.Dockerfile,
				Hostname:   tt.fields.Hostname,
				Aliases:    tt.fields.Aliases,
				Path:       tt.fields.Path,
				Webroot:    tt.fields.Webroot,
				Extensions: tt.fields.Extensions,
				Xdebug:     tt.fields.Xdebug,
				Blackfire:  tt.fields.Blackfire,
				Database:   tt.fields.Database,
				HomeDir:    tt.fields.HomeDir,
			}
			if got := a.GetHostname(); got != tt.want {
				t.Errorf("GetHostname() = %v, want %v", got, tt.want)
			}
		})
	}
}
