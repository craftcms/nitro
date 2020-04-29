package nitro

import (
	"runtime"

	"github.com/craftcms/nitro/validate"
)

func SSH(name string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}

	syscall := true
	if runtime.GOOS == "windows" {
		syscall = false
	}

	return &Action{
		Type:       "shell",
		UseSyscall: syscall,
		Args:       []string{"shell", name},
	}, nil
}
