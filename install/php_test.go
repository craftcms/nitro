package install

import (
	"strings"
	"testing"
)

func TestPHPReturnsCommands(t *testing.T) {
	// Arrange
	v := "7.4"
	expected := 13

	// Act
	commands, err := PHP(v)
	total := len(strings.Split(commands, " "))

	// Assert
	if err != nil {
		t.Errorf("expected the error to be nil, got %v", err.Error())
	}
	if total != expected {
		t.Errorf("expected total number of commands to be %v, got %v instead", expected, total)
	}
}

func TestPHPReturnsError(t *testing.T) {
	// Arrange
	v := "nothere"

	// Act
	_, err := PHP(v)

	// Assert
	if err == nil {
		t.Errorf("expected the error to NOT be nil, got %v", err.Error())
	}
}
