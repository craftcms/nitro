package nitro

import (
	"fmt"

	"github.com/craftcms/nitro/validate"
)

func UnmountDir(name, target string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}

	return &Action{
		Type:       "umount",
		Output:     fmt.Sprintf("Unmounting %s from %s", target, name),
		UseSyscall: false,
		Args:       []string{"umount", name + ":" + target},
	}, nil
}
