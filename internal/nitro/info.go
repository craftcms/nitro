package nitro

func Info(name string) []Command {
	return []Command{
		{
			Machine: name,
			Type:    "info",
			Args:    []string{name},
		},
	}
}
