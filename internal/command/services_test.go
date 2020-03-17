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
	set.Bool("mysql", false, "docs")
	c := cli.NewContext(nil, set, nil)
	r := &TestRunner{}

	type args struct {
		c *cli.Context
		r *TestRunner
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
			expectedArgs: []string{"exec", "nitro-test", "--", "sudo", "service", "nginx", "restart"},
			toParse:      []string{"--nginx"},
			wantErr:      false,
		},
		// TODO fix this later
		//{
		//	name:         "restart mysql",
		//	args:         args{c: c, r: r},
		//	expectedArgs: []string{"exec", "nitro-test", "--", "sudo", "service", "mariadb", "start"},
		//	toParse:      []string{"--mysql"},
		//	wantErr:      false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = set.Parse(tt.toParse)

			if err := serviceRestartAction(tt.args.c, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("serviceRestartAction() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.args.r.args) == 0 {
				t.Errorf("expected args to not be zero; got %q instead", tt.args.r.args)
			}

			for i, arg := range tt.args.r.args {
				if tt.expectedArgs[i] != arg {
					t.Errorf("expected the arg %q; got %q instead", tt.expectedArgs[i], arg)
					t.Log(tt.args.r.args)
				}
			}
		})
	}
}
