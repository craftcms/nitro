package nitro

import (
	"errors"
	"fmt"
	"strings"

	"github.com/craftcms/nitro/validate"
)

func NginxReload(name string) (*Action, error) {
	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "sudo", "service", "nginx", "restart"},
	}, nil
}

func CreateSiteSymllink(name, site string) (*Action, error) {
	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "sudo", "ln", "-s", "/etc/nginx/sites-available/" + site, "/etc/nginx/sites-enabled/"},
	}, nil
}

func CopyNginxTemplate(name, hostname string) (*Action, error) {
	if hostname == "" {
		return nil, errors.New("hostname cannot be empty")
	}
	if err := validate.Hostname(hostname); err != nil {
		return nil, err
	}

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "sudo", "cp", "/opt/nitro/nginx/template.conf", "/etc/nginx/sites-available/" + hostname},
	}, nil
}

func CreateNginxSiteDirectory(name, site string) (*Action, error) {
	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "mkdir", "-p", "/nitro/sites/" + site},
	}, nil
}

func ChangeTemplateVariables(name, webroot, hostname, php string, aliases []string) (*[]Action, error) {
	var actions []Action
	template := hostname
	if aliases != nil {
		hostname = hostname + " " + strings.Join(aliases, " ")
	}

	actions = append(actions, *ChangeNginxTemplateVariable(name, template, "CHANGEWEBROOTDIR", webroot))
	actions = append(actions, *ChangeNginxTemplateVariable(name, template, "CHANGESERVERNAME", hostname))
	actions = append(actions, *ChangeNginxTemplateVariable(name, template, "CHANGEPHPVERSION", php))

	return &actions, nil
}

func ChangeNginxTemplateVariable(machine, hostname, variable, actual string) *Action {
	sedCmd := "s|" + variable + "|" + actual + "|g"

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", machine, "--", "sudo", "sed", "-i", sedCmd, fmt.Sprintf("/etc/nginx/sites-available/%v", hostname)},
	}
}
