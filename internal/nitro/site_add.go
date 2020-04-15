package nitro

import "fmt"

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

func CopyNginxTemplate(name, site string) (*Action, error) {
	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "sudo", "cp", "/opt/nitro/nginx/template.conf", "/etc/nginx/sites-available/" + site},
	}, nil
}

func CreateNginxSiteDirectory(name, site string) (*Action, error) {
	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "mkdir", "-p", "/nitro/sites/" + site},
	}, nil
}

func ChangeTemplateVariables(name, site, dir, php string) (*[]Action, error) {
	var actions []Action

	actions = append(actions, *changeVariables(name, site, "CHANGEPATH", site))
	actions = append(actions, *changeVariables(name, site, "CHANGESERVERNAME", site))
	actions = append(actions, *changeVariables(name, site, "CHANGEPUBLICDIR", dir))
	actions = append(actions, *changeVariables(name, site, "CHANGEPHPVERSION", php))

	return &actions, nil
}

func changeVariables(name, site, variable, actual string) *Action {
	file := fmt.Sprintf("/etc/nginx/sites-available/%v", site)

	sedCmd := "s|" + variable + "|" + actual + "|g"

	return &Action{
		Type:       "exec",
		UseSyscall: false,
		Args:       []string{"exec", name, "--", "sudo", "sed", "-i", sedCmd, file},
	}
}
