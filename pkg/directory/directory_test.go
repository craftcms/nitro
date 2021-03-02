package directory

import (
	"path/filepath"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "non empty directories return false",
			args: args{dir: filepath.Join("testdata")},
			want: false,
		},
		{
			name: "missing directories return false",
			args: args{dir: filepath.Join("missing")},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.args.dir); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}
