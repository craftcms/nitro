package command

import (
	"flag"
	"testing"

	"github.com/urfave/cli/v2"
)

func Test_serviceRestartAction(t *testing.T) {
	// Arrange
	set := flag.NewFlagSet("test", 0)
	set.String("machine", "nitro-test", "doc")
	set.Bool("nginx", false, "docs")
	c := cli.NewContext(nil, set, nil)
	r := &TestRunner{}

	type args struct {
		c *cli.Context
		r Runner
	}
	tests := []struct {
		name         string
		args         args
		expectedArgs []string
		toParse      []string
		wantErr      bool
	}{
		{
			name:         "restart nginx",
			args:         args{c: c, r: r},
			expectedArgs: []string{"exec", "nitro-text", "--", "sudo", "service", "nginx", "restart"},
			toParse:      []string{"--nginx"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = set.Parse(tt.toParse)

			if err := serviceRestartAction(tt.args.c, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("serviceRestartAction() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
