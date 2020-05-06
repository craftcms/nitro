package nitro

import (
	"fmt"

	"github.com/craftcms/nitro/validate"
)

func RemoveSymlink(name, site string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}
	if err := validate.Hostname(site); err != nil {
		return nil, err
	}

	return &Action{
		Type:       "exec",
		Output:     fmt.Sprintf("Removing symlink for %s", site),
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "sudo", "rm", "/etc/nginx/sites-enabled/" + site},
	}, nil
}

func RemoveNginxSiteDirectory(name, site string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}
	if err := validate.Hostname(site); err != nil {
		return nil, err
	}

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "rm", "-rf", "/app/sites/" + site},
	}, nil
}
