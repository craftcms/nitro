package command

import (
	"flag"
	"testing"

	"github.com/urfave/cli/v2"
)

type TestRunner struct {
	path    string
	Args    []string
	syscall bool
	input   string
}

func (r *TestRunner) Run(args []string) error {
	r.Args = args

	return nil
}

func (r TestRunner) UseSyscall(t bool) {
	r.syscall = t
}

func (r TestRunner) Path() string {
	return r.path
}

func (r TestRunner) SetInput(input string) error {
	r.input = input
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

func Test_removeAction(t *testing.T) {
	// Arrange
	set := flag.NewFlagSet("test", 0)
	set.String("machine", "nitro-test", "doc")
	c := cli.NewContext(nil, set, nil)
	_ = set.Parse([]string{"--machine=test", "somehost"})
	expectedArgs := []string{"exec", "test", "--", "sudo", "bash", "/opt/nitro/nginx/remove-site.sh", "somehost"}
	r := TestRunner{}

	// Act
	err := removeAction(c, &r)

	// Assert
	if err != nil {
		t.Error("expected error to be nil")
	}
	if assertArgsMatch(expectedArgs, r.Args) == false {
		t.Errorf("expected the Args to match; got %v instead\n", r.Args)
	}
}

func assertArgsMatch(expected []string, actual []string) bool {
	if len(expected) != len(actual) {
		return false
	}
	for i, v := range expected {
		if v != actual[i] {
			return false
		}
	}

	return true
}
