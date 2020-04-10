package action

import "errors"

func Redis(name string) (*Action, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "redis-cli"},
	}, nil
}
