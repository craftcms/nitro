package bindmounts

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFromDir(t *testing.T) {
	// get the cwd
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// create the path base for the tests
	base := filepath.Join(wd, "testdata")

	type args struct {
		path     string
		excludes []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "returns the correct number of bind mounts from a directory with no excludes",
			args: args{
				path:     filepath.Join(base, "project-with-composer-deps"),
				excludes: []string{"node_modules"},
			},
			want: []string{
				filepath.Join(base, "project-with-composer-deps", "app"),
				filepath.Join(base, "project-with-composer-deps", "config"),
				filepath.Join(base, "project-with-composer-deps", "vendor"),
			},
		},
		{
			name: "returns the correct number of bind mounts from a directory with no excludes",
			args: args{
				path:     filepath.Join(base, "project-with-composer-deps"),
				excludes: []string{"node_modules", "vendor"},
			},
			want: []string{
				filepath.Join(base, "project-with-composer-deps", "app"),
				filepath.Join(base, "project-with-composer-deps", "config"),
			},
		},
		{
			name: "no excludes returns all directories",
			args: args{
				path: filepath.Join(base, "project-with-composer-deps"),
			},
			want: []string{
				filepath.Join(base, "project-with-composer-deps", "app"),
				filepath.Join(base, "project-with-composer-deps", "config"),
				filepath.Join(base, "project-with-composer-deps", "node_modules"),
				filepath.Join(base, "project-with-composer-deps", "vendor"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromDir(tt.args.path, tt.args.excludes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromDir() = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}
