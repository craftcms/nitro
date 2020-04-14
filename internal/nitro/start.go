package nitro

import "errors"

func Start(name string) (*Action, error) {
	if name == "" {
		return nil, errors.New("machine name cannot be empty")
	}

	return &Action{
		Type:       "start",
		UseSyscall: false,
		Args:       []string{"start", name},
	}, nil
}
