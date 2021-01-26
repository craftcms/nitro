package filetype

import (
	"path/filepath"
	"testing"
)

func TestDetermine(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "tarfile.tar.gz returns application/zip",
			args: args{
				file: filepath.Join("testdata", "tarfile.tar.gz"),
			},
			want:    "tar",
			wantErr: false,
		},
		{
			name: "example.zip returns application/zip",
			args: args{
				file: filepath.Join("testdata", "example.zip"),
			},
			want:    "zip",
			wantErr: false,
		},
		{
			name: "backup.sql returns",
			args: args{
				file: filepath.Join("testdata", "backup.sql"),
			},
			want:    "text",
			wantErr: false,
		},
		{
			name: "example.txt returns text/plain",
			args: args{
				file: filepath.Join("testdata", "example.txt"),
			},
			want:    "text",
			wantErr: false,
		},
		{
			name: "directory returns error",
			args: args{
				file: filepath.Join("testdata"),
			},
			wantErr: true,
		},
		{
			name: "no file returns error",
			args: args{
				file: "nowhere",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Determine(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Determine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Determine() = %v, want %v", got, tt.want)
			}
		})
	}
}
