package nitro

import "errors"

// Info will display the machine information based on a name
func Info(name string) (*Action, error) {
	if name == "" {
		return nil, errors.New("missing machine name")
	}

	return &Action{
		Type:       "info",
		UseSyscall: false,
		Args:       []string{"info", name},
	}, nil
}
