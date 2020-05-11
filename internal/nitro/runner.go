package nitro

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type ShellRunner interface {
	Path() string
	UseSyscall(t bool)
	SetInput(input string) error
	Run(args []string) error
}

type MultipassRunner struct {
	path       string
	useSyscall bool
	stdin      *io.Reader
	input      string
}

func (m *MultipassRunner) Path() string {
	return m.path
}

func (m *MultipassRunner) UseSyscall(t bool) {
	m.useSyscall = t
}

func (m *MultipassRunner) SetInput(input string) error {
	if input == "" {
		return errors.New("input must not be empty")
	}

	// set the input
	m.input = input

	return nil
}

func (m *MultipassRunner) Run(args []string) error {
	// if this is a syscall, hand it off
	if m.useSyscall {
		// this allows us to only add Args, not the binary path to
		// keep everything consistent
		if args[0] != "multipass" {
			args = append([]string{"multipass"}, args...)
		}

		return syscall.Exec(m.path, args, os.Environ())
	}

	cmd := exec.Command(m.path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if m.input != "" {
		cmd.Stdin = strings.NewReader(m.input)
	}

	return cmd.Run()
}

func NewMultipassRunner(file string) ShellRunner {
	path, err := exec.LookPath(file)
	if err != nil {
		log.Fatal("Unable to find multipass")
	}

	return &MultipassRunner{
		path: path,
	}
}
