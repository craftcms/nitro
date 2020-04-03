package nitro

func Redis(name string) []Command {
	return []Command{
		{
			Machine: name,
			Type:    "exec",
			Args:    []string{name, "--", "redis-cli"},
		},
	}
}
