package v3

import (
	"os"
	"path/filepath"
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
			want: "custom-hostname-from-file.nitro",
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
			want: "mysite.nitro",
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
