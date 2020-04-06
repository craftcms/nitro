package nitro

func Info(name string) []Command {
	return []Command{
		{
			Machine:   name,
			Chainable: true,
			Type:      "info",
			Args:      []string{name},
		},
	}
}
