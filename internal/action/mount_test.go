package action

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
			name: "valid args returns action",
			args: args{
				name:   "somename",
				folder: "/tmp",
				site:   "example.test",
			},
			want: &Action{
				Type:       "mount",
				UseSyscall: false,
				Args:       []string{"mount", "/tmp", "somename:/app/sites/example.test"},
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
				site:   "example",
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
				t.Errorf("Mount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMountDirectory(t *testing.T) {
	type args struct {
		name        string
		source      string
		destination string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "valid args returns action",
			args: args{
				name:        "somename",
				source:      "/tmp",
				destination: "/home/ubuntu/sites/test",
			},
			want: &Action{
				Type:       "mount",
				UseSyscall: false,
				Args:       []string{"mount", "/tmp", "somename:/home/ubuntu/sites/test"},
			},
			wantErr: false,
		},
		{
			name: "invalid source returns error",
			args: args{
				name:        "somename",
				source:      "not-here",
				destination: "/home/ubuntu/dev",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid name returns error",
			args: args{
				name:        "",
				source:      "/tmp",
				destination: "/home/ubuntu/dev",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MountDirectory(tt.args.name, tt.args.source, tt.args.destination)
			if (err != nil) != tt.wantErr {
				t.Errorf("MountDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MountDirectory() got = %v, want %v", got, tt.want)
			}
		})
	}
}
