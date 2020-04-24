package task

import (
	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

func Rename(machine, php string, site config.Site) ([]nitro.Action, error) {
	var actions []nitro.Action

	// remove the symlink
	removeSymlinkAction, err := nitro.RemoveSymlink(machine, site.Hostname)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *removeSymlinkAction)

	copyTemplateAction, err := nitro.CopyNginxTemplate(machine, site.Hostname)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *copyTemplateAction)

	changeNginxVariablesAction, err := nitro.ChangeTemplateVariables(machine, site.Webroot, site.Hostname, php, site.Aliases)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *changeNginxVariablesAction...)

	// restart nginx
	restartNginxAction, err := nitro.NginxReload(machine)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *restartNginxAction)

	return actions, nil
}
