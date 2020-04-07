package nitro

import "fmt"

func CreateNewDirectoryForSite(name, site string) Command {
	return Command{
		Machine:   name,
		Type:      "exec",
		Chainable: true,
		Args:      []string{name, "--", "mkdir", "-p", "/app/sites/" + site},
	}
}

func CopyNginxTemplate(name, site string) Command {
	return Command{
		Machine:   name,
		Type:      "exec",
		Chainable: true,
		Args:      []string{name, "--", "sudo", "cp", "/opt/nitro/nginx/template.conf", "/etc/nginx/sites-available/" + site},
	}
}

func ChangeVariablesInTemplate(name, domain, dir, php string) []Command {
	var commands []Command

	commands = append(commands, changeVariableInNginxTemplate(name, domain, "CHANGEPATH", domain))
	commands = append(commands, changeVariableInNginxTemplate(name, domain, "CHANGESERVERNAME", domain))
	commands = append(commands, changeVariableInNginxTemplate(name, domain, "CHANGEPUBLICDIR", dir))
	commands = append(commands, changeVariableInNginxTemplate(name, domain, "CHANGEPHPVERSION", php))

	return commands
}

func changeVariableInNginxTemplate(name, site, variable, actual string) Command {
	file := fmt.Sprintf("/etc/nginx/sites-available/%v", site)
	sedCmd := "s|" + variable + "|" + actual + "|g"
	return Command{
		Machine:   name,
		Type:      "exec",
		Chainable: true,
		Args:      []string{name, "--", "sudo", "sed", "-i", sedCmd, file},
	}
}

func LinkNginxSite(name, site string) Command {
	return Command{
		Machine:   name,
		Type:      "exec",
		Chainable: true,
		Args:      []string{name, "--", "sudo", "ln", "-s", "/etc/nginx/sites-available/" + site, "/etc/nginx/sites-enabled/"},
	}
}

func ReloadNginx(name string) Command {
	return Command{
		Machine:   name,
		Type:      "exec",
		Chainable: true,
		Args:      []string{name, "--", "sudo", "service", "nginx", "restart"},
	}
}
