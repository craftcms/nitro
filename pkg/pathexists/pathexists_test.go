package pathexists

import (
	"path/filepath"
	"testing"
)

func TestIsDirectory(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "existing directory returns true",
			args: args{
				path: filepath.Join("testdata", "exists"),
			},
			want: true,
		},
		{
			name: "existing files returns false",
			args: args{
				path: filepath.Join("testdata", "exists", ".gitigore"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDirectory(tt.args.path); got != tt.want {
				t.Errorf("IsDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "existing file returns true",
			args: args{
				path: filepath.Join("testdata", "exists", ".gitignore"),
			},
			want: true,
		},
		{
			name: "existing directory returns false",
			args: args{
				path: filepath.Join("testdata", "exists"),
			},
			want: false,
		},
		{
			name: "non-existing file returns false",
			args: args{
				path: filepath.Join("testdata", "exists", "missing"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFile(tt.args.path); got != tt.want {
				t.Errorf("IsFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
