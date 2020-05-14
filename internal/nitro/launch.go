package nitro

import (
	"errors"
	"strconv"

	"github.com/craftcms/nitro/validate"
)

type Action struct {
	Type       string
	UseSyscall bool
	Input      string
	Args       []string
}

// Launch is responsible for the creation of a virtual machine, each parameter must be provided and validated
// prior to making the machine. The input param needs to be a valid cloud-config string.
func Launch(name string, cpus int, memory, disk, input string) (*Action, error) {
	if name == "" {
		return nil, errors.New("the name of the machine cannot be empty")
	}
	if cpus == 0 {
		return nil, errors.New("the number of CPUs cannot be 0")
	}
	if err := validate.DiskSize(disk); err != nil {
		return nil, err
	}
	if err := validate.Memory(memory); err != nil {
		return nil, err
	}
	if input == "" {
		return nil, errors.New("input cannot be empty")
	}

	return &Action{
		Type:       "launch",
		UseSyscall: false,
		Input:      input,
		Args:       []string{"launch", "--name", name, "--cpus", strconv.Itoa(cpus), "--mem", memory, "--disk", disk, "18.04", "--cloud-init", "-"},
	}, nil
}
