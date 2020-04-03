package nitro

import "fmt"

func RestartDatabase(name, engine, version string) []Command {
	imageName := fmt.Sprintf("nitro_%v_%v", engine, version)
	fmt.Println(imageName)
	return []Command{
		{
			Machine: name,
			Type:    "exec",
			Args:    []string{name, "--", "docker", "restart", imageName},
		},
	}
}
