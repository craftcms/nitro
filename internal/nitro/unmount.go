package nitro

import (
	"github.com/craftcms/nitro/validate"
)

func UnmountDir(name, target string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}

	return &Action{
		Type:       "umount",
		UseSyscall: false,
		Args:       []string{"umount", name + ":" + target},
	}, nil
}
