package scripts

import (
	"reflect"
	"testing"
)

func TestAptInstall(t *testing.T) {
	type args struct {
		name string
		pkgs []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "",
			args: args{
				name: "some-machine-name",
				pkgs: []string{"some", "packages"},
			},
			want: []string{"exec", "some-machine-name", "--", "sudo", "apt", "install", "-y", "some", "packages"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AptInstall(tt.args.name, tt.args.pkgs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AptInstall() = %v, want %v", got, tt.want)
			}
		})
	}
}
