package nitro

import (
	"github.com/craftcms/nitro/validate"
)

func RestartPhpFpm(name, php string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}
	if err := validate.PHPVersion(php); err != nil {
		return nil, err
	}

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "sudo", "service", "php" + php + "-fpm", "restart"},
	}, nil
}
