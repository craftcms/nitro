package appaware

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/craftcms/nitro/pkg/config"
)

func TestDetect(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testdir := filepath.Join(wd, "testdata")

	type args struct {
		cfg config.Config
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "returns an error when there is no app",
			args: args{
				cfg: config.Config{
					ParsedApps: []config.App{
						{
							Hostname:   "existing-app.nitro",
							Path:       "~/existing-app",
							PHPVersion: "8.0",
						},
					},
					HomeDirectory: testdir,
				},
				dir: filepath.Join(testdir, "not-an-existing-app", "web"),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "can find an existing app based on the sub-directory",
			args: args{
				cfg: config.Config{
					Apps: []config.App{
						{
							Hostname:   "existing-app.nitro",
							Path:       "~/existing-app",
							PHPVersion: "8.0",
						},
					},
					HomeDirectory: testdir,
				},
				dir: filepath.Join(testdir, "existing-app", "web"),
			},
			want:    "existing-app.nitro",
			wantErr: false,
		},
		{
			name: "can find an existing app based on the directory",
			args: args{
				cfg: config.Config{
					Apps: []config.App{
						{
							Hostname:   "existing-app.nitro",
							Path:       "~/existing-app",
							PHPVersion: "8.0",
						},
					},
					HomeDirectory: testdir,
				},
				dir: filepath.Join(testdir, "existing-app"),
			},
			want:    "existing-app.nitro",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Detect(tt.args.cfg, tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("Detect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Detect() \ngot = \n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}
