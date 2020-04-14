package nitro

import (
	"errors"

	"github.com/craftcms/nitro/validate"
)

func EnableXdebug(name, php string) (*Action, error) {
	if name == "" {
		return nil, errors.New("name cannot by empty")
	}
	if err := validate.PHPVersion(php); err != nil {
		return nil, err
	}

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "sudo", "phpenmod", "-v", php, "xdebug"},
	}, nil
}

func DisableXdebug(name, php string) (*Action, error) {
	if name == "" {
		return nil, errors.New("name cannot by empty")
	}
	if err := validate.PHPVersion(php); err != nil {
		return nil, err
	}

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "sudo", "phpdismod", "-v", php, "xdebug"},
	}, nil
}
