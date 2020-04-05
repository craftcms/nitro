package scripts

import "errors"

func DockerExec(args []string) (*Script, error) {
	if len(args) > 0 {
		return &Script{}, errors.New("no arguments where provided to docker exec")
	}

	cmd := []string{"docker", "exec"}
	cmd = append(cmd, args...)

	return &Script{
		Name: "running docker exec ",
		Args: cmd,
	}, nil
}
