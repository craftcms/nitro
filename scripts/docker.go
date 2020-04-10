package scripts

import (
	"fmt"
)

// DockerRunDatabase builds the commands to run a database inside of docker for a machine.
func DockerRunDatabase(name, engine, version string) []string {
	var port string
	var envvars []string
	switch engine {
	case "postgres":
		port = "5432"
		envvars = []string{"-e", "POSTGRES_PASSWORD=nitro", "-e", "POSTGRES_USER=nitro", "-e", "POSTGRES_DB=nitro"}
	default:
		port = "3306"
		envvars = []string{"-e", "MYSQL_ROOT_PASSWORD=nitro", "-e", "MYSQL_DATABASE=nitro", "-e", "MYSQL_USER=nitro", "-e", "MYSQL_PASSWORD=nitro"}
	}

	// TODO clean this up
	imageName := fmt.Sprintf("nitro_%v_%v", engine, version)
	hostPath := fmt.Sprintf("/opt/nitro/volumes/%v", engine)
	var containerPath string
	if engine == "mysql" {
		containerPath = "/var/lib/mysql"
	} else {
		containerPath = "/var/lib/postgresql/data"
	}
	volumeMount := fmt.Sprintf("%v:%v", hostPath, containerPath)

	commands := []string{name, "--", "docker", "run", "-v", volumeMount, "--name", imageName, "-d", "--restart=always", "-p", port + ":" + port}

	// append the environment variables
	commands = append(commands, envvars...)

	// append the docker image as the last arg
	image := []string{engine + ":" + version}
	commands = append(commands, image...)

	return commands
}
