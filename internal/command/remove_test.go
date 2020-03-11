package command

import (
	"flag"
	"fmt"
	"testing"

	"github.com/urfave/cli/v2"
)

type TestRunner struct {
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

func (r TestRunner) SetInput(input string) error {
	return nil
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
