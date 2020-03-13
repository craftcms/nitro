package command

import (
	"flag"
	"testing"

	"github.com/urfave/cli/v2"
)

type TestRunner struct {
	path    string
	args    []string
	syscall bool
	input   string
}

func (r TestRunner) Run(args []string) error {
	r.args = args

	return nil
}

func (r TestRunner) UseSyscall(t bool) {
	r.syscall = t
}

func (r TestRunner) Path() string {
	return r.path
}

func (r TestRunner) SetInput(input string) error {
	return nil
}

func Test_removeBeforeAction(t *testing.T) {
	// Arrange
	set := flag.NewFlagSet("test", 0)
	set.String("machine", "nitro-test", "doc")
	c := cli.NewContext(nil, set, nil)

	type args struct {
		c *cli.Context
	}
	tests := []struct {
		name    string
		args    args
		toParse []string
		wantErr bool
	}{
		{
			name:    "requires a host argument",
			args:    args{c},
			toParse: []string{"--machine=nitro-test"},
			wantErr: true,
		},
		{
			name:    "no error when a host argument is provided",
			args:    args{c},
			toParse: []string{"--machine=nitro-test", "hostname"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = set.Parse(tt.toParse)
			if err := removeBeforeAction(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("removeBeforeAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
