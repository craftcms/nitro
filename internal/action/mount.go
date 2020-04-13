package action

import "github.com/craftcms/nitro/validate"

func Mount(name, folder, site string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}
	if err := validate.Path(folder); err != nil {
		return nil, err
	}
	if err := validate.Domain(site); err != nil {
		return nil, err
	}

	return &Action{
		Type:       "mount",
		UseSyscall: false,
		Args:       []string{"mount", folder, name + ":/app/sites/" + site},
	}, nil
}

func MountDirectory(name, source, destination string) (*Action, error) {
	return &Action{
		Type:       "mount",
		UseSyscall: false,
		Args:       []string{"mount", source, name + ":" + destination},
	}, nil
}
