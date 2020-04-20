package helpers

import (
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestNormalizePath(t *testing.T) {
	home, err := homedir.Dir()
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	projectDir := strings.Replace(currentDir, home, "~", 1)

	type args struct {
		path string
		home string
	}
	tests := []struct {
		name            string
		args            args
		homedir         string
		wantFilename    string
		wantFileAbsPath string
		wantErr         bool
	}{
		{
			name:            "will resolve as a full path",
			args:            args{path: "testdata/normalize-path/somefile.txt", home: home},
			wantFilename:    "somefile.txt",
			wantFileAbsPath: currentDir + "/testdata/normalize-path/somefile.txt",
			wantErr:         false,
		},
		{
			name:            "will resolve ~ as a full path",
			args:            args{path: projectDir + "/testdata/normalize-path/somefile.txt", home: home},
			wantFilename:    "somefile.txt",
			wantFileAbsPath: currentDir + "/testdata/normalize-path/somefile.txt",
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := NormalizePath(tt.args.path, tt.args.home)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantFilename {
				t.Errorf("NormalizePath() got = \n%v, \nwant \n%v", got, tt.wantFilename)
			}
			if got1 != tt.wantFileAbsPath {
				t.Errorf("NormalizePath() got1 = \n%v, \nwant \n%v", got1, tt.wantFileAbsPath)
			}
		})
	}
}
