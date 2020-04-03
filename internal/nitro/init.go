package nitro

import (
	"github.com/craftcms/nitro/scripts"
)

const cloudConfig = `#cloud-config
packages:
  - redis
  - jq
  - apt-transport-https
  - ca-certificates
  - curl
  - gnupg-agent
  - software-properties-common
runcmd:
  - sudo add-apt-repository -y ppa:nginx/stable
  - sudo add-apt-repository -y ppa:ondrej/php
  - curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
  - sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
  - sudo apt-get update -y
  - sudo apt install -y nginx docker-ce docker-ce-cli containerd.io
  - sudo usermod -aG docker ubuntu
  - wget -q -O - https://packages.blackfire.io/gpg.key | sudo apt-key add -
  - echo "deb http://packages.blackfire.io/debian any main" | sudo tee /etc/apt/sources.list.d/blackfire.list
`

type Command struct {
	Machine   string
	Type      string
	Chainable bool
	Input     string
	Args      []string
}

func Init(name, cpus, memory, disk, php, db, version string) []Command {
	var commands []Command

	// add the init command
	commands = append(commands, Command{
		Machine:   name,
		Type:      "launch",
		Chainable: true,
		Input:     cloudConfig,
		Args:      []string{"--name", name, "--cpus", cpus, "--mem", memory, "--disk", disk, "--cloud-init", "-"},
	})

	// install the core packages
	installCommands := []string{name, "--", "sudo", "apt", "install", "-y"}
	installCommands = append(installCommands, scripts.InstallPHP(php)...)
	commands = append(commands, Command{
		Machine:   name,
		Chainable: true,
		Type:      "exec",
		Args:      installCommands,
	})

	var port string
	var envvars []string
	switch db {
	case "postgres":
		port = "5432"
		envvars = []string{"-e", "POSTGRES_PASSWORD=nitro", "-e", "POSTGRES_USER=nitro", "-e", "POSTGRES_DB=nitro"}
	default:
		port = "3306"
		envvars = []string{"-e", "MYSQL_ROOT_PASSWORD=nitro", "-e", "MYSQL_DATABASE=nitro", "-e", "MYSQL_USER=nitro", "-e", "MYSQL_PASSWORD=nitro"}
	}

	// setup the docker commands
	dockerCommands := []string{name, "--", "docker", "run", "-d", "--restart=always", "-p", port + ":" + port}
	dockerCommands = append(dockerCommands, envvars...)
	image := []string{db + ":" + version}
	dockerCommands = append(dockerCommands, image...)
	commands = append(commands, Command{
		Machine:   name,
		Chainable: true,
		Type:      "exec",
		Args:      dockerCommands,
	})

	return commands
}
