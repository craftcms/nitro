package nitro

import (
	"reflect"
	"testing"
)

func TestSQL(t *testing.T) {
	type args struct {
		name    string
		engine  string
		version string
		root    bool
	}
	tests := []struct {
		name string
		args args
		want []Command
	}{
		{
			name: "",
			args: args{
				name:    "",
				engine:  "",
				version: "",
				root:    false,
			},
			want: []Command{
				{

				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SQL(tt.args.name, tt.args.engine, tt.args.version, tt.args.root); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
