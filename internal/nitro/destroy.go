package nitro

func Destroy(name string, permanent bool) []Command {
	var commands []Command
	if permanent {
		commands = append(commands, Command{
			Machine:   name,
			Chainable: false,
			Type:      "delete",
			Args:      []string{name, "-p"},
		})
	} else {
		commands = append(commands, Command{
			Machine:   name,
			Chainable: false,
			Type:      "delete",
			Args:      []string{name},
		})
	}

	return commands
}
