package nitro

import (
	"reflect"
	"testing"
)

func TestRemoveSymlink(t *testing.T) {
	type args struct {
		name string
		site string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "can remove valid symlink",
			args: args{
				name: "somemachine",
				site: "example.test",
			},
			want: &Action{
				Type:       "exec",
				Output:     "Removing symlink for example.test",
				UseSyscall: false,
				Args:       []string{"exec", "somemachine", "--", "sudo", "rm", "/etc/nginx/sites-enabled/example.test"},
			},
			wantErr: false,
		},
		{
			name: "invalid name returns error",
			args: args{
				name: "",
				site: "example.test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid site returns error",
			args: args{
				name: "somename",
				site: "not valid",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RemoveSymlink(tt.args.name, tt.args.site)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveSymlink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveSymlink() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveNginxSiteDirectory(t *testing.T) {
	type args struct {
		name string
		site string
	}
	tests := []struct {
		name    string
		args    args
		want    *Action
		wantErr bool
	}{
		{
			name: "can remove site directory",
			args: args{
				name: "somemachine",
				site: "example.test",
			},
			want: &Action{
				Type:       "exec",
				UseSyscall: false,
				Args:       []string{"exec", "somemachine", "--", "rm", "-rf", "/app/sites/example.test"},
			},
			wantErr: false,
		},
		{
			name: "invalid name returns error",
			args: args{
				name: "",
				site: "example.test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid site returns error",
			args: args{
				name: "somename",
				site: "not valid",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RemoveNginxSiteDirectory(tt.args.name, tt.args.site)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveNginxSiteDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveNginxSiteDirectory() got = %v, want %v", got, tt.want)
			}
		})
	}
}
