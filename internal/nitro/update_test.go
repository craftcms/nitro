package nitro

import (
	"reflect"
	"testing"
)

func TestUpdate(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want []Command
	}{
		{
			name: "gets the right args",
			args: args{name: "this"},
			want: []Command{
				{
					Machine: "this",
					Type:    "exec",
					Args:    []string{"this", "--", "sudo", "apt-get", "upgrade", "-y"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Update(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
