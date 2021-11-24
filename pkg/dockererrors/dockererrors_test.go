package dockererrors

import (
	"errors"
	"testing"
)

func TestIsPortError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "returns true if the error is talking about a port collision",
			args: args{err: errors.New("port is already allocated")},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPortError(tt.args.err); got != tt.want {
				t.Errorf("IsPortError() = %v, want %v", got, tt.want)
			}
		})
	}
}
