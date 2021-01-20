package checkfile

import (
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "returns true when the file exists",
			args: args{
				path: filepath.Join("testdata", "with-file", "example.txt"),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "returns false when the file does not exist",
			args: args{
				path: filepath.Join("testdata", "no-directory", "example.txt"),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Exists(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}
