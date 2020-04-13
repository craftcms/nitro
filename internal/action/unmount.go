package action

import "github.com/craftcms/nitro/validate"

func Unmount(name, site string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}

	return &Action{
		Type:       "umount",
		UseSyscall: false,
		Args:       []string{"umount", name + ":/app/sites/" + site},
	}, nil
}
