package nitro

func Start(name string) []Command {
	return []Command{
		{
			Machine: name,
			Type:    "start",
			Args:    []string{name},
		},
	}
}
