package nitrod

import "os/exec"

// Runner is an interface to run commands.
type Runner interface {
	Run(command string, args []string) ([]byte, error)
}

// ServiceRunner is an implementation of the Runner interface
// that uses the exec.Command
type ServiceRunner struct{}

// Run sends the commands provided to exec.Command
func (r ServiceRunner) Run(command string, args []string) ([]byte, error) {
	return exec.Command(command, args...).CombinedOutput()
}
