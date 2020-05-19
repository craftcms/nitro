package nitro

import (
	"reflect"
	"testing"
)

func TestMount(t *testing.T) {
	type args struct {
		name   string
		folder string
		site   string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "valid args returns nitro",
			args: args{
				name:   "somename",
				folder: "/tmp",
				site:   "example.test",
			},
			want: &Action{
				Type:       "mount",
				UseSyscall: false,
				Args:       []string{"mount", "/tmp", "somename:/home/ubuntu/sites/example.test"},
			},
			wantErr: false,
		},
		{
			name: "invalid name returns error",
			args: args{
				name:   "",
				folder: "somefolder",
				site:   "example.test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid path returns error",
			args: args{
				name:   "somename",
				folder: "not-here",
				site:   "example.test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid domain returns error",
			args: args{
				name:   "somename",
				folder: "/tmp",
				site:   "example test",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mount(tt.args.name, tt.args.folder, tt.args.site)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mount() got = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
}
