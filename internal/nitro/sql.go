package nitro

func SQL(name, engine, version string, root bool) []Command {
	var dockerArgs []string
	switch engine {
	// TODO add the postgres args
	default:
		if root {
			dockerArgs = append(dockerArgs, "mysql", "-u", "root", "-pnitro")
		} else {
			dockerArgs = append(dockerArgs, "mysql", "-u", "nitro", "-pnitro")
		}
	}

	args := []string{name, "--", "docker", "exec", "-it", "nitro_" + engine + "_" + version}
	args = append(args, dockerArgs...)

	return []Command{
		{
			Machine:   name,
			Type:      "exec",
			Chainable: false,
			Args:      args,
		},
	}
}
