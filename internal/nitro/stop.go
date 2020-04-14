package nitro

import "errors"

func Stop(name string) (*Action, error) {
	if name == "" {
		return nil, errors.New("machine name cannot be empty")
	}

	return &Action{
		Type:       "stop",
		UseSyscall: false,
		Args:       []string{"stop", name},
	}, nil
}
