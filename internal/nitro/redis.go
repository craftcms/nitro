package nitro

import (
	"errors"
	"runtime"
)

func Redis(name string) (*Action, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	syscall := true
	if runtime.GOOS == "windows" {
		syscall = false
	}

	return &Action{
		Type:       "exec",
		UseSyscall: syscall,
		Args:       []string{"exec", name, "--", "redis-cli"},
	}, nil
}
