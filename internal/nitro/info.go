package nitro

import "errors"

// Info will display the machine information based on a name
func Info(machine string) (*Action, error) {
	if machine == "" {
		return nil, errors.New("missing machine name")
	}

	return &Action{
		Type:       "info",
		UseSyscall: false,
		Args:       []string{"info", machine},
	}, nil
}
