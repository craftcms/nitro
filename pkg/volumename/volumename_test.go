package volumename

import (
	"os"
	"strings"
	"testing"
)

func TestFromPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "complete example to catch them all",
			args: args{
				path: strings.Join([]string{"this", "that", "OR", "anot:her"}, string(os.PathSeparator)),
			},
			want: "this_that_or_anot_her",
		},
		{
			name: "upper case paths are transformed to lowercase",
			args: args{
				path: "UPPER",
			},
			want: "upper",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromPath(tt.args.path); got != tt.want {
				t.Errorf("FromPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
