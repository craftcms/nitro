package action

import "github.com/craftcms/nitro/validate"

func RemoveSymlink(name, site string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "sudo", "rm", "/etc/nginx/sites-enabled/" + site},
	}, nil
}

func RemoveNginxSiteDirectory(name, site string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "rm", "-rf", "/app/sites/" + site},
	}, nil
}
