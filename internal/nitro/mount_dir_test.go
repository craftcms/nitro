package nitro

import (
	"reflect"
	"testing"
)

func TestMountDir(t *testing.T) {
	type args struct {
		name   string
		source string
		target string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "mounts a directory",
			args: args{
				name:   "somename",
				source: "./testdata/source-folder",
				target: "home/ubuntu/sites",
			},
			want: &Action{
				Type:       "mount",
				UseSyscall: false,
				Args:       []string{"mount", "./testdata/source-folder", "somename:/home/ubuntu/sites"},
			},
			wantErr: false,
		},
		{
			name: "mounts a directory and removes the trailing slash",
			args: args{
				name:   "somename",
				source: "./testdata/source-folder/",
				target: "/home/ubuntu/sites",
			},
			want: &Action{
				Type:       "mount",
				UseSyscall: false,
				Args:       []string{"mount", "./testdata/source-folder", "somename:/home/ubuntu/sites"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MountDir(tt.args.name, tt.args.source, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("MountDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MountDir() got = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}
