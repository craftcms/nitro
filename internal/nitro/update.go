package nitro

func SSH(name string) []Command {
	return []Command{
		{
			Machine: name,
			Type:    "shell",
			Args:    []string{"shell", name},
		},
	}
}
