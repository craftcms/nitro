package nitro

import (
	"reflect"
	"testing"
)

func TestDestroy(t *testing.T) {
	type args struct {
		name      string
		permanent bool
	}
	tests := []struct {
		name string
		args args
		want []Command
	}{
		{
			name: "get command to permanently delete the machine",
			args: args{
				name:      "some-machine",
				permanent: false,
			},
			want: []Command{
				{
					Machine:   "some-machine",
					Chainable: false,
					Type:      "delete",
					Args:      []string{"some-machine"},
				},
			},
		},
		{
			name: "get command to permanently delete the machine",
			args: args{
				name:      "some-machine",
				permanent: true,
			},
			want: []Command{
				{
					Machine:   "some-machine",
					Chainable: false,
					Type:      "delete",
					Args:      []string{"some-machine", "-p"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Destroy(tt.args.name, tt.args.permanent); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Destroy() = %v, want %v", got, tt.want)
			}
		})
	}
}
