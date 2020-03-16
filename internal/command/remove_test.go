package command

import (
	"flag"
	"fmt"
	"testing"

	"github.com/urfave/cli/v2"
)

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
