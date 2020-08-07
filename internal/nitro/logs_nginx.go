package nitro

import "github.com/craftcms/nitro/internal/validate"

func LogsNginx(name, kind string) (*Action, error) {
	if err := validate.MachineName(name); err != nil {
		return nil, err
	}

	switch kind {
	case "access":
		return &Action{
			Type:       "exec",
			UseSyscall: false,
			Args:       []string{"exec", name, "--", "sudo", "tail", "-f", "/var/log/nginx/access.log"},
		}, nil
	case "error":
		return &Action{
			Type:       "exec",
			UseSyscall: false,
			Args:       []string{"exec", name, "--", "sudo", "tail", "-f", "/var/log/nginx/error.log"},
		}, nil
	}

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "sudo", "tail", "-f", "/var/log/nginx/access.log", "-f", "/var/log/nginx/error.log"},
	}, nil
}
