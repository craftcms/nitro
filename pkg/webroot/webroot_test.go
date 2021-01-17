package webroot

import (
	"path/filepath"
	"testing"
)

func TestFind(t *testing.T) {
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
			name: "find the webroot",
			args: args{
				path: filepath.Join("testdata", "no-vendor"),
			},
			want:    "web",
			wantErr: false,
		},
		{
			name: "vendor paths are ignored",
			args: args{
				path: filepath.Join("testdata", "with-vendor"),
			},
			want:    "public",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Find(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Find() = %v, want %v", got, tt.want)
			}
		})
	}
}
