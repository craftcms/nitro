package task

import (
	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

func Add(machine string, configFile config.Config, site config.Site) ([]nitro.Action, error) {
	var actions []nitro.Action
	// mount the directory
	m := configFile.Mounts[len(configFile.Mounts)-1]
	mountAction, err := nitro.MountDir(machine, m.AbsSourcePath(), m.Dest)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *mountAction)

	// copy the nginx template
	copyTemplateAction, err := nitro.CopyNginxTemplate(machine, site.Hostname)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *copyTemplateAction)

	// copy the nginx template
	changeNginxVariablesAction, err := nitro.ChangeTemplateVariables(machine, site.Webroot, site.Hostname, configFile.PHP, site.Aliases)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *changeNginxVariablesAction...)

	createSymlinkAction, err := nitro.CreateSiteSymllink(machine, site.Hostname)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *createSymlinkAction)

	restartNginxAction, err := nitro.NginxReload(machine)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *restartNginxAction)

	return actions, nil
}