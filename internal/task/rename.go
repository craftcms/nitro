package task

import (
	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

func RenameSite(machine, php string, oldSite, newSite config.Site) ([]nitro.Action, error) {
	var actions []nitro.Action

	// remove the symlink from the old oldSite
	removeSymlinkAction, err := nitro.RemoveSymlink(machine, oldSite.Hostname)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *removeSymlinkAction)

	// create a new oldSite config for the new hostname
	copyTemplateAction, err := nitro.CopyNginxTemplate(machine, newSite.Hostname)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *copyTemplateAction)

	// change the webroot path
	actions = append(
		actions,
		*nitro.ChangeNginxTemplateVariable(machine, newSite.Hostname, "CHANGEWEBROOTDIR", newSite.Webroot),
	)

	// change the server name variable
	actions = append(
		actions,
		*nitro.ChangeNginxTemplateVariable(machine, newSite.Hostname, "CHANGESERVERNAME", newSite.Hostname),
	)

	// change the PHP version
	actions = append(
		actions,
		*nitro.ChangeNginxTemplateVariable(machine, newSite.Hostname, "CHANGEPHPVERSION", php),
	)

	// reload nginx
	reloadNginxAction, err := nitro.NginxReload(machine)
	if err != nil {
		return nil, err
	}

	actions = append(actions, *reloadNginxAction)

	return actions, nil
}

func Rename(machine, php string, existingSite, renamedSite config.Site, mount *config.Mount) ([]nitro.Action, error) {
	var actions []nitro.Action

	// remove the symlink from the old oldSite
	removeSymlinkAction, err := nitro.RemoveSymlink(machine, existingSite.Hostname)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *removeSymlinkAction)

	// create a new oldSite config
	copyTemplateAction, err := nitro.CopyNginxTemplate(machine, renamedSite.Hostname)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *copyTemplateAction)

	changeNginxVariablesAction, err := nitro.ChangeTemplateVariables(machine, existingSite.Webroot, existingSite.Hostname, php, existingSite.Aliases)
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

	if mount != nil {
		/// unmount the directory
		unMountAction, err := nitro.Unmount(machine, existingSite.Hostname)
		if err != nil {
			return nil, err
		}
		actions = append(actions, *unMountAction)

		// mount the new directory
		mountAction, err := nitro.Mount(machine, mount.AbsSourcePath(), renamedSite.Hostname)
		if err != nil {
			return nil, err
		}
		actions = append(actions, *mountAction)
	}

	return actions, nil
}
