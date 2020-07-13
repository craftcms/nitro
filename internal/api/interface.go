package api

import "os/exec"

type Runner interface {
	Run(command string, args []string) ([]byte, error)
}

type ServiceRunner struct {}

func (r ServiceRunner) Run(command string, args []string) ([]byte, error) {
	return exec.Command(command, args...).CombinedOutput()
}
