package nitro

import (
	"strings"

	"github.com/craftcms/nitro/internal/validate"
)

func UnmountDir(name, target string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}

	if strings.Contains(target, " ") {
		target = strings.TrimSpace(target)
	}

	return &Action{
		Type:       "umount",
		UseSyscall: false,
		Args:       []string{"umount", name + ":" + target},
	}, nil
}
