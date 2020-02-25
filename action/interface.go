package action

import (
	"log"
	"os/exec"
	"syscall"
)

// CommandLineExecutor is an interface that handles
// running commands, the real use case is to call syscall.Exec
type CommandLineExecutor interface {
	// Path will use exec.LookPath to find a complete path to a file, in real world use, this would be
	// a path to the multipass binary
	Path() string

	// Exec matches syscall.Exec args so we can perform assertions without directly calling the
	// underlying system.
	Exec(argv0 string, argv []string, envv []string) (err error)
}

// MockExecutor is used for successful mocking
type MockExecutor struct {
	path      string
	argv0     string
	arguments []string
	env       []string
}

func (m *MockExecutor) Path() string {
	return "testing"
}

func (s *MockExecutor) Exec(argv0 string, argv []string, envv []string) (err error) {
	s.argv0 = argv0
	s.arguments = argv
	s.env = envv
	return nil
}

type SyscallExecutor struct {
	// the path to the executable file
	path string
}

func (s *SyscallExecutor) Path() string {
	return s.path
}

func (s *SyscallExecutor) Exec(argv0 string, argv []string, envv []string) (err error) {
	return syscall.Exec(argv0, argv, envv)
}

// NewSyscallExecutor will lookup a file path and
// return a new SyscallExecutor struct with that path
func NewSyscallExecutor(file string) *SyscallExecutor {
	path, err := exec.LookPath(file)
	if err != nil {
		log.Fatal(err)
	}

	return &SyscallExecutor{path: path}
}
