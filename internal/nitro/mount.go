package nitro

import (
	"github.com/craftcms/nitro/internal/validate"
)

// Mount is used to mount a folder for a specific website
func Mount(name, folder, site string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}
	if err := validate.Path(folder); err != nil {
		return nil, err
	}
	if err := validate.Hostname(site); err != nil {
		return nil, err
	}

	return &Action{
		Type:       "mount",
		UseSyscall: false,
		Args:       []string{"mount", folder, name + ":" + "/home/ubuntu/sites/" + site},
	}, nil
}
