package v3

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestConfig_GetAppHostName(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	homeDir := filepath.Join(wd, "testdata")

	type fields struct {
		Containers []Container
		Blackfire  Blackfire
		Databases  []Database
		Services   Services
		Apps       []App
		HomeDir    string
		ConfigFile string
	}
	type args struct {
		hostname string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "apps without a top-level config keys take precedence",
			args: args{hostname: "mysite.nitro"},
			fields: fields{
				Apps: []App{
					{
						Config: "~/path-with-config/nitro.yaml",
					},
				},
				HomeDir: homeDir,
			},
			want:    "custom-hostname-from-file.nitro",
			wantErr: false,
		},
		{
			name: "apps with top-level config keys take precedence",
			args: args{hostname: "mysite.nitro"},
			fields: fields{
				Apps: []App{
					{
						Config:   "~/path-with-config/nitro.yaml",
						Hostname: "mysite.nitro",
					},
				},
				HomeDir: homeDir,
			},
			want:    "mysite.nitro",
			wantErr: false,
		},
		{
			name: "apps with no hostname return an error",
			args: args{hostname: "missing.nitro"},
			fields: fields{
				Apps: []App{
					{
						PHPVersion: "8.0",
					},
				},
				HomeDir: homeDir,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{
				Containers: tt.fields.Containers,
				Blackfire:  tt.fields.Blackfire,
				Databases:  tt.fields.Databases,
				Services:   tt.fields.Services,
				Apps:       tt.fields.Apps,
				HomeDir:    tt.fields.HomeDir,
				ConfigFile: tt.fields.ConfigFile,
			}

			got, err := c.GetAppHostName(tt.args.hostname)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAppHostName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAppHostName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

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
						PHPVersion: "8.0",
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
