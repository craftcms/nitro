package nitro

import "fmt"

func NginxLogs(name, kind string) []Command {
	var commands []Command
	switch kind {
	case "access":
		commands = append(commands, Command{
			Machine: name,
			Type:    "exec",
			Args:    []string{name, "--", "sudo", "tail", "-f", "/var/log/nginx/access.log"},
		})
	case "error":
		commands = append(commands, Command{
			Machine: name,
			Type:    "exec",
			Args:    []string{name, "--", "sudo", "tail", "-f", "/var/log/nginx/error.log"},
		})
	default:
		commands = append(commands, Command{
			Machine: name,
			Type:    "exec",
			Args:    []string{name, "--", "sudo", "tail", "-f", "/var/log/nginx/access.log", "-f", "/var/log/nginx/error.log"},
		})
	}

	return commands
}

func DatabaseLogs(name, engine, version string) []Command {
	imageName := fmt.Sprintf("nitro_%v_%v", engine, version)
	return []Command{
		{
			Machine: name,
			Type:    "exec",
			Args:    []string{name, "--", "docker", "logs", "-f", imageName},
		},
	}
}
