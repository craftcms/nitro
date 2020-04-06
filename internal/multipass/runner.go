package multipass

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// Multipasser is an interface that provides a set of helpful wrappers around executing multipass
// commands. It provides the option to pass input (e.g. when using the multipass launch command
// to pass a cloud config) and a method to tell the command to use syscall (when chaining commands
// is not required for execution (e.g. multipass shell commands)
type Multipasser interface {
	UseSyscall()
	SetInput(input string) error
	Run(args []string) error
}

type multipass struct {
	Syscall bool
	Input   string
	Path    string
}

func (m *multipass) UseSyscall() {
	m.Syscall = true
}

func (m *multipass) SetInput(input string) error {
	if input == "" {
		return errors.New("input cannot be empty")
	}

	m.Input = input

	return nil
}

func (m *multipass) Run(args []string) error {
	if m.Path == "" {
		return errors.New("path to multipass is not defined")
	}

	if m.Syscall {
		return syscall.Exec(m.Path, args, os.Environ())
	}

	// create a new command
	c := exec.Command(m.Path, args...)

	// if we have input, pass it to stdin for the command
	if m.Input != "" {
		c.Stdin = strings.NewReader(m.Input)
	}

	if c.Stdout == nil {
		c.Stdout = os.Stdout
	}

	if c.Stderr == nil {
		c.Stderr = os.Stderr
	}

	return c.Run()
}

// New takes a file path (e.g. multipass) and will look for the full path and
// return a new Multipasser or an error. This provides a simple setup to
// use multipass.
func New(file string) (Multipasser, error) {
	path, err := exec.LookPath(file)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &multipass{
		Path: path,
	}, nil
}
