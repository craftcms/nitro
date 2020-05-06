package nitro

import (
	"errors"
	"fmt"
)

func Restart(machine string) (*Action, error) {
	if machine == "" {
		return nil, errors.New("machine cannot be empty")
	}

	return &Action{
		Type:       "restart",
		Output:     fmt.Sprintf("Restarting machine %q", machine),
		UseSyscall: false,
		Args:       []string{"restart", machine},
	}, nil
}
