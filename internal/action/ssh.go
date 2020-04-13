package action

import "github.com/craftcms/nitro/validate"

func SSH(name string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}

	return &Action{
		Type:       "shell",
		UseSyscall: true,
		Args:       []string{"shell", name},
	}, nil
}
