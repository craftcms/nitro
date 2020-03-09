package command

import (
	"testing"

	"github.com/urfave/cli/v2"
)

func Test_removeBeforeAction(t *testing.T) {
	type args struct {
		c *cli.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "returns error when no args are provided",
			wantErr: true,
		},
		{
			name:    "does not return error when args are provided",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := removeBeforeAction(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("removeBeforeAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
