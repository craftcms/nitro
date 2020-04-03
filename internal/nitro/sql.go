package nitro

func SQL(name, engine, version string) []Command {
	return []Command{
		{
			Machine:   name,
			Type:      "exec",
			Chainable: false,
			Args:      []string{name, "--", "docker", "exec", "-it", "nitro_" + engine + "_" + version, "mysql"},
		},
	}
}
