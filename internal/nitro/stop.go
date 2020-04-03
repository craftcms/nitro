package nitro

func Stop(name string) []Command {
	return []Command{
		{
			Machine: name,
			Type:    "stop",
			Args:    []string{name},
		},
	}
}
