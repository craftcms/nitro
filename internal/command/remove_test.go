package command

import (
	"flag"
	"fmt"
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

func TestContextFlag(t *testing.T) {
	set := flag.NewFlagSet("test", 0)
	set.Bool("myflag", false, "doc")
	c := cli.NewContext(nil, set, nil)
	_ = set.Parse([]string{"--myflag", "bat", "baz"})

	if c.Args().Len() != 2 {
		t.Error("length should be two")
	}

	if c.Bool("myflag") == false {
		t.Error("myflag should be true")
	}
}

func TestRemoveBeforeAction(t *testing.T) {
	// Arrange
	set := flag.NewFlagSet("test", 0)
	set.String("machine", "nitro-test", "doc")
	c := cli.NewContext(nil, set, nil)
	_ = set.Parse([]string{"--machine=nitro-test"})

	// act
	err := removeBeforeAction(c)

	// Assert
	if err != ErrRemoveNoHostArgProvided {
		t.Errorf("expected error to be %v, got %v instead", "no host was specified for removal", err)
	}
}

func TestRemoveBeforeCommandReturnsError(t *testing.T) {
	// Arrange
	set := flag.NewFlagSet("test", 0)
	ctx := cli.NewContext(nil, set, nil)
	expected := "no host was specified for removal"

	// Act
	err := removeBeforeAction(ctx)

	// Assert
	if err == nil {
		t.Error("expected the error from removeBeforeAction() to not be nil")
	}
	if err.Error() != expected {
		t.Errorf("expected the error from removeBeforeAction() to be %v; got %v instead", expected, err)
	}
	fmt.Println(ctx.Context.Value("host"))
}

func TestRemoveAfterActionPrintsOutput(t *testing.T) {
	// Arrange
	set := flag.NewFlagSet("test", 0)
	ctx := cli.NewContext(nil, set, nil)

	// Act
	if err := removeAfterAction(ctx); err != nil {
		t.Errorf("expected error to be nil. got %v instead", err)
	}
}
