package nitro

func Update(name string) []Command {
	return []Command{
		{
			Machine: name,
			Type:    "exec",
			Args:    []string{name, "--", "sudo", "apt-get", "upgrade", "-y"},
		},
	}
}
