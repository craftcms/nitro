package task

import (
	"fmt"

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
		exists, _ := configFile.AlreadyMounted(mount)
		if !exists {
			unmountAction, err := nitro.UnmountDir(machine, mount.Dest)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *unmountAction)
			fmt.Println("Removing mount", mount.Source, "from", machine)

			actions = append(actions, nitro.Action{
				Type:       "exec",
				UseSyscall: false,
				Args:       []string{"exec", machine, "--", "rm", "-rf", mount.Dest},
			})
		}
	}

	// check if there are mounts we need to create
	for _, mount := range configFile.Mounts {
		exists, _ := inMemoryConfig.AlreadyMounted(mount)
		if !exists {
			mountAction, err := nitro.MountDir(machine, mount.AbsSourcePath(), mount.Dest)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *mountAction)
			fmt.Println("Mounting", mount.Source, "to", machine)
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

			// remove symlink
			removeSiteAvailable, err := nitro.RemoveNginxSiteAvailable(machine, site.Hostname)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *removeSiteAvailable)

			// reload nginx
			reloadNginxAction, err := nitro.NginxReload(machine)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *reloadNginxAction)
			fmt.Println("Removing site", site.Hostname, "from", machine)
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

			createSymlink, err := nitro.CreateSiteSymllink(machine, site.Hostname)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *createSymlink)

			// reload nginx
			reloadNginxAction, err := nitro.NginxReload(machine)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *reloadNginxAction)
			fmt.Println("Adding site", site.Hostname, "to", machine)
		}
	}

	// check if there are databases to remove
	for _, database := range inMemoryConfig.Databases {
		if !configFile.DatabaseExists(database) {
			actions = append(actions, nitro.Action{
				Type:       "exec",
				UseSyscall: false,
				Args:       []string{"exec", machine, "--", "docker", "rm", "-v", database.Name(), "-f"},
			})
			fmt.Println("Removing database", database.Name(), "from", machine)
		}
	}

	// check if there are database to create
	for _, database := range configFile.Databases {
		if !inMemoryConfig.DatabaseExists(database) {
			createVolume, err := nitro.CreateDatabaseVolume(machine, database.Engine, database.Version, database.Port)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *createVolume)

			createContainer, err := nitro.CreateDatabaseContainer(machine, database.Engine, database.Version, database.Port)
			if err != nil {
				return nil, err
			}
			actions = append(actions, *createContainer)

			fmt.Println("Creating database", database.Name(), "on", machine)
		}
	}

	// if the php versions do not match, install the requested version - which makes it the default
	if configFile.PHP != php {
		// install the php version
		installPhp, err := nitro.InstallPackages(machine, configFile.PHP)
		if err != nil {
			return nil, err
		}
		actions = append(actions, *installPhp)

		// set the default php
		setPhpDefault := &nitro.Action{
			Type:       "exec",
			UseSyscall: false,
			Args:       []string{"exec", machine, "--", "sudo", "update-alternatives", "--set", "php", "/usr/bin/php" + configFile.PHP},
		}
		actions = append(actions, *setPhpDefault)

		// set the default phar
		setDefaultPhar := &nitro.Action{
			Type:       "exec",
			UseSyscall: false,
			Args:       []string{"exec", machine, "--", "sudo", "update-alternatives", "--set", "phar", "/usr/bin/phar" + configFile.PHP},
		}
		actions = append(actions, *setDefaultPhar)

		// set the default phar.phar
		setDefaultPharPhar := &nitro.Action{
			Type:       "exec",
			UseSyscall: false,
			Args:       []string{"exec", machine, "--", "sudo", "update-alternatives", "--set", "phar.phar", "/usr/bin/phar.phar" + configFile.PHP},
		}
		actions = append(actions, *setDefaultPharPhar)

		// set the default phpize
		setDefaultPhpize := &nitro.Action{
			Type:       "exec",
			UseSyscall: false,
			Args:       []string{"exec", machine, "--", "sudo", "update-alternatives", "--set", "phpize", "/usr/bin/phpize" + configFile.PHP},
		}
		actions = append(actions, *setDefaultPhpize)

		// set the default php-config
		setDefaultPhpConfig := &nitro.Action{
			Type:       "exec",
			UseSyscall: false,
			Args:       []string{"exec", machine, "--", "sudo", "update-alternatives", "--set", "php-config", "/usr/bin/php-config" + configFile.PHP},
		}
		actions = append(actions, *setDefaultPhpConfig)

		fmt.Println("Installing PHP", configFile.PHP, "on", machine)
	}

	return actions, nil
}
