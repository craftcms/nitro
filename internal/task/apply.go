package task

import (
	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

// Apply is responsible for comparing the current configuration and what information is
// found on a machine such as fromMultipassMounts and sites. Apple will then take the appropriate
// steps to compare are create actions that "normal up" the configuration state.
func Apply(machine string, configFile config.Config, mounts []config.Mount, sites []config.Site, dbs []config.Database, php string) ([]nitro.Action, error) {
	var actions []nitro.Action
	inMemoryConfig := config.Config{PHP: php, Mounts: mounts, Sites: sites, Databases: dbs}

	// check if there are mounts we need to remove
	for _, mount := range inMemoryConfig.Mounts {
		if !configFile.MountExists(mount.Dest) {
			unmountAction, err := nitro.UnmountDir(machine, mount.Dest)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *unmountAction)
		}
	}

	// check if there are mounts we need to create
	for _, mount := range configFile.Mounts {
		if !inMemoryConfig.MountExists(mount.Dest) {
			mountAction, err := nitro.MountDir(machine, mount.AbsSourcePath(), mount.Dest)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *mountAction)
		}
	}

	// check if there are sites we need to remove
	for _, site := range inMemoryConfig.Sites {
		if !configFile.SiteExists(site) {
			// remove symlink
			removeSymlink, err := nitro.RemoveSymlink(machine, site.Hostname)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *removeSymlink)

			// reload nginx
			reloadNginxAction, err := nitro.NginxReload(machine)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *reloadNginxAction)
		}
	}

	// check if there are sites we need to make
	for _, site := range configFile.Sites {
		// find the parent to mount
		if !inMemoryConfig.SiteExists(site) {
			// copy template
			copyTemplateAction, err := nitro.CopyNginxTemplate(machine, site.Hostname)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *copyTemplateAction)

			// replace variable
			changeNginxVariablesAction, err := nitro.ChangeTemplateVariables(machine, site.Webroot, site.Hostname, configFile.PHP, site.Aliases)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *changeNginxVariablesAction...)

			// reload nginx
			reloadNginxAction, err := nitro.NginxReload(machine)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *reloadNginxAction)
		}
	}

	return actions, nil
}
