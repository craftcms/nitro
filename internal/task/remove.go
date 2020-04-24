package task

import (
	"github.com/craftcms/nitro/config"
	"github.com/craftcms/nitro/internal/nitro"
)

func Remove(name string, mount config.Mount, site config.Site) ([]nitro.Action, error) {
	var actions []nitro.Action

	// unmount
	unmountAction, err := nitro.UnmountDir(name, mount.Dest)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *unmountAction)

	// remove nginx symlink
	removeSymlinkAction, err := nitro.RemoveSymlink(name, site.Hostname)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *removeSymlinkAction)

	restartNginxAction, err := nitro.NginxReload(name)
	if err != nil {
		return nil, err
	}
	actions = append(actions, *restartNginxAction)

	return actions, nil
}
