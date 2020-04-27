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

	for _, mount := range inMemoryConfig.Mounts {
		if !configFile.MountExists(mount.Dest) {
			unmountAction, err := nitro.UnmountDir(machine, mount.Dest)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *unmountAction)
		}
	}

	for _, mount := range configFile.Mounts {
		if !inMemoryConfig.MountExists(mount.Dest) {
			mountAction, err := nitro.MountDir(machine, mount.AbsSourcePath(), mount.Dest)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *mountAction)
		}
	}

	return actions, nil
}
