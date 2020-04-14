package nitro

import "errors"

func Restart(name string) (*Action, error) {
	if name == "" {
		return nil, errors.New("machine name cannot be empty")
	}

	return &Action{
		Type:       "restart",
		UseSyscall: false,
		Args:       []string{"restart", name},
	}, nil
}
