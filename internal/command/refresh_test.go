package command

import (
	"flag"
	"testing"

	"github.com/urfave/cli/v2"
)

func Test_refreshAction(t *testing.T) {
	// Arrange
	set := flag.NewFlagSet("test", 0)
	set.String("machine", "nitro-test", "doc")
	c := cli.NewContext(nil, set, nil)
	r := SpyTestRunner{}

	type args struct {
		c *cli.Context
		r SpyTestRunner
	}
	tests := []struct {
		name         string
		args         args
		expectedArgs []string
		toParse      []string
		wantErr      bool
	}{
		{
			name:    "",
			args:    args{c: c, r: r},
			toParse: []string{"refresh"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = set.Parse(tt.toParse)

			if err := refreshAction(tt.args.c, &tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("refreshAction() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Log(tt.args.r.args)
		})
	}
}
