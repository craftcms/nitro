package nitro

import (
	"errors"
)

func Restart(machine string) (*Action, error) {
	if machine == "" {
		return nil, errors.New("machine cannot be empty")
	}

	return &Action{
		Type:       "restart",
		UseSyscall: false,
		Args:       []string{"restart", machine},
	}, nil
}
