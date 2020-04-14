package nitro

import (
	"errors"

	"github.com/craftcms/nitro/validate"
)

func LogsDocker(name, container string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}
	if container == "" {
		return nil, errors.New("container name is empty")
	}

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "docker", "logs", container, "-f"},
	}, nil
}
