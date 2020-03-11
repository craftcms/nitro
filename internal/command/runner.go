package command

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// CmdRunner is a struct that will handle calling exec.Cmd or syscall.Exec using a boolean
// it keeps all of the logic in a single interface to keep testing easier
type CmdRunner struct {
	path      string
	isSyscall bool
	stdin     *io.Reader
	input     string
}

func (c *CmdRunner) SetInput(input string) error {
	if input == "" {
		return errors.New("input should not be an empty string")
	}

	c.input = input

	return nil
}

func (c CmdRunner) Run(args []string) error {
	// if this is a syscall, hand it off
	if c.isSyscall {
		// this allows us to only add args, not the binary path to
		// keep everything consistent
		if args[0] != "multipass" {
			args = append([]string{"multipass"}, args...)
		}

		return syscall.Exec(c.path, args, os.Environ())
	}

	cmd := exec.Command(c.path, args...)

	if c.input != "" {
		cmd.Stdin = strings.NewReader(c.input)
	}

	if cmd.Stdout == nil {
		cmd.Stdout = os.Stdout
	}

	if cmd.Stderr == nil {
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}

func (c *CmdRunner) UseSyscall(t bool) {
	c.isSyscall = t
}

func (c CmdRunner) Path() string {
	return c.path
}

func NewRunner(file string) Runner {
	path, err := exec.LookPath(file)
	if err != nil {
		log.Fatal("unable to find multipass")
	}

	return &CmdRunner{
		path: path,
	}
}
