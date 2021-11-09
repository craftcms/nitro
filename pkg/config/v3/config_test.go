package v3

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	home := filepath.Join(wd, "testdata")

	type args struct {
		home string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "Can load a config and get the correct values based with overrides",
			args: args{home: home},
			want: &Config{
				Apps: []App{
					{
						Hostname:   "mysite.nitro",
						Path:       "~/mysite",
						PHPVersion: "7.4",
						Webroot:    "web",
					},
					{
						Config:     "~/team-site-with-config/nitro.yaml",
						Hostname:   "team-site-name-from-config.nitro",
						Aliases:    []string{"my-local-app.test"},
						PHPVersion: "8.0",
						Extensions: []string{"grpc"},
						Webroot:    "public",
					},
				},
				HomeDir: home,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.home)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got:\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}
