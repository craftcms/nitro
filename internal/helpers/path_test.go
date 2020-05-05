package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParentPathName(t *testing.T) {
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
			name: "can get the path directory name",
			args: args{
				path: "./testdata/good-example",
			},
			want:    "good-example",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PathName(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("PathName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PathName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDirectoryArg(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	abspath, err := filepath.Abs(wd)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		args []string
	}
	tests := []struct {
		name         string
		args         args
		directory    string
		absolutePath string
		wantErr      bool
	}{
		{
			name:         "can get the directory and abs path when not sending args",
			args:         args{args: nil},
			directory:    "helpers",
			absolutePath: abspath,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetDirectoryArg(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDirectoryArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.directory {
				t.Errorf("GetDirectoryArg() got = %v, want %v", got, tt.directory)
			}
			if got1 != tt.absolutePath {
				t.Errorf("GetDirectoryArg() got1 = %v, want %v", got1, tt.absolutePath)
			}
		})
	}
}
