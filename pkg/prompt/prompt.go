package prompt

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/phpversions"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/pkg/validate"
	"github.com/craftcms/nitro/pkg/webroot"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// CreateDatabase is used to interactively walk a user through creating a new database. It will return true if the user created a database along
// with the hostname, database, port, and driver for the database container.
func CreateDatabase(cmd *cobra.Command, docker client.CommonAPIClient, output terminal.Outputer) (bool, string, string, string, string, error) {
	confirm, err := output.Confirm("Add a database for the site?", true, "")
	if err != nil {
		return false, "", "", "", "", err
	}

	if !confirm {
		return false, "", "", "", "", nil
	}

	// make sure the context is not nil
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// add filters to show only the environment and database containers
	filter := filters.NewArgs()
	filter.Add("label", containerlabels.Nitro)
	filter.Add("label", containerlabels.Type+"=database")

	// get a list of all the databases
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter, All: true})
	if err != nil {
		return false, "", "", "", "", err
	}

	// sort containers by the name
	sort.SliceStable(containers, func(i, j int) bool {
		return containers[i].Names[0] < containers[j].Names[0]
	})

	// get all of the containers as a list
	var engineOpts []string
	for _, c := range containers {
		// start the container if not running
		if c.State != "running" {
			for _, command := range cmd.Root().Commands() {
				if command.Use == "start" {
					if err := command.RunE(cmd, []string{}); err != nil {
						return false, "", "", "", "", err
					}
				}
			}
		}

		engineOpts = append(engineOpts, strings.TrimLeft(c.Names[0], "/"))
	}

	// prompt the user for the engine to add the database
	var containerID, databaseEngine string
	selected, err := output.Select(os.Stdin, "Select the database engine: ", engineOpts)
	if err != nil {
		return false, "", "", "", "", err
	}

	// set the container id and db engine
	containerID = containers[selected].ID
	databaseEngine = containers[selected].Labels[containerlabels.DatabaseCompatibility]
	if containerID == "" {
		return false, "", "", "", "", fmt.Errorf("unable to get the container")
	}

	// ask the user for the database to create
	db, err := output.Ask("Enter the new database name", "", ":", &validate.DatabaseName{})
	if err != nil {
		return false, "", "", "", "", err
	}

	output.Pending("creating database", db)

	// set the commands based on the engine type
	var cmds, privileges []string
	switch databaseEngine {
	case "mysql":
		cmds = []string{"mysql", "-uroot", "-pnitro", fmt.Sprintf(`-e CREATE DATABASE IF NOT EXISTS %s;`, db)}
		privileges = []string{"mysql", "-uroot", "-pnitro", fmt.Sprintf(`-e GRANT ALL PRIVILEGES ON * TO '%s'@'%s';`, "nitro", "%")}
	default:
		cmds = []string{"psql", "--username=nitro", "--host=127.0.0.1", fmt.Sprintf(`-c CREATE DATABASE %s;`, db)}
	}

	// create the exec
	e, err := docker.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd:          cmds,
	})
	if err != nil {
		return false, "", "", "", "", err
	}

	// attach to the container
	resp, err := docker.ContainerExecAttach(ctx, e.ID, types.ExecStartCheck{
		Tty: false,
	})
	if err != nil {
		return false, "", "", "", "", err
	}
	defer resp.Close()

	// start the exec
	if err := docker.ContainerExecStart(ctx, e.ID, types.ExecStartCheck{}); err != nil {
		return false, "", "", "", "", fmt.Errorf("unable to start the container, %w", err)
	}

	// check if we should grant privileges
	if privileges != nil {
		// create the exec
		e, err := docker.ContainerExecCreate(ctx, containerID, types.ExecConfig{
			AttachStdout: true,
			AttachStderr: true,
			Tty:          false,
			Cmd:          privileges,
		})
		if err != nil {
			return false, "", "", "", "", err
		}

		// attach to the container
		resp, err := docker.ContainerExecAttach(ctx, e.ID, types.ExecStartCheck{
			Tty: false,
		})
		if err != nil {
			return false, "", "", "", "", err
		}
		defer resp.Close()

		// start the exec
		if err := docker.ContainerExecStart(ctx, e.ID, types.ExecStartCheck{}); err != nil {
			return false, "", "", "", "", fmt.Errorf("unable to start the container, %w", err)
		}

		// wait for the container exec to complete
		waiting := true
		for waiting {
			resp, err := docker.ContainerExecInspect(ctx, e.ID)
			if err != nil {
				return false, "", "", "", "", err
			}

			waiting = resp.Running
		}
	}

	output.Done()

	output.Info("Database added 💪")

	// get the container hostname
	hostname := strings.TrimLeft(containers[selected].Names[0], "/")

	// get the info from the container
	info, err := docker.ContainerInspect(ctx, containers[selected].ID)
	if err != nil {
		return false, "", "", "", "", err
	}

	var port string
	for p := range info.NetworkSettings.Ports {
		if port != "" {
			break
		}

		port = p.Port()
	}

	// set the driver for the database
	driver := "mysql"
	if containers[selected].Labels[containerlabels.DatabaseCompatibility] == "postgres" {
		driver = "pgsql"
	}

	return true, hostname, db, port, driver, nil
}

// CreateSite takes the users home directory and the site path and walked the user
// through adding a site to the config.
func CreateSite(home, dir string, output terminal.Outputer) (*config.Site, error) {
	// create a new site
	site := config.Site{}

	// get the hostname from the directory
	// p := filepath.Join(dir)
	sp := strings.Split(filepath.Join(dir), string(os.PathSeparator))
	site.Hostname = sp[len(sp)-1]

	// append the test domain if there are no periods
	if !strings.Contains(site.Hostname, ".") {
		// set the default tld
		tld := "nitro"
		if os.Getenv("NITRO_DEFAULT_TLD") != "" {
			tld = os.Getenv("NITRO_DEFAULT_TLD")
		}

		site.Hostname = fmt.Sprintf("%s.%s", site.Hostname, tld)
	}

	// prompt for the hostname
	hostname, err := output.Ask("Enter the hostname", site.Hostname, ":", &validate.HostnameValidator{})
	if err != nil {
		return nil, err
	}

	// set the input as the hostname
	site.Hostname = hostname

	output.Success("setting hostname to", site.Hostname)

	// set the sites directory but make the path relative
	siteAbsPath, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	site.Path = strings.Replace(siteAbsPath, home, "~", 1)

	output.Success("adding site", site.Path)

	// get the web directory
	found, _ := webroot.Find(dir)

	// if the root is still empty, we fall back to the default
	if found == "" {
		found = "web"
	}

	// set the web root
	site.Webroot = found

	// prompt for the web root
	root, err := output.Ask("Enter the web root for the site", site.Webroot, ":", nil)
	if err != nil {
		return nil, err
	}

	site.Webroot = root

	output.Success("using web root", site.Webroot)

	// prompt for the php version
	versions := phpversions.Versions
	selected, err := output.Select(os.Stdin, "Choose a PHP version: ", versions)
	if err != nil {
		return nil, err
	}

	// set the version of php
	site.Version = versions[selected]

	output.Success("setting PHP version", site.Version)

	// load the config
	cfg, err := config.Load(home)
	if err != nil {
		return nil, err
	}

	// add the site to the config
	if err := cfg.AddSite(site); err != nil {
		return nil, err
	}

	// save the config file
	if err := cfg.Save(); err != nil {
		return nil, err
	}

	return &site, nil
}

// RunApply will prompt a user to run the apply command. It optionally accepts a "force"
// option that will not prompt the user and run apply regardless.
func RunApply(cmd *cobra.Command, args []string, force bool, output terminal.Outputer) error {
	if !force {
		// ask if the apply command should run
		apply, err := output.Confirm("Apply changes now?", true, "")
		if err != nil {
			return err
		}

		// if apply is false return nil
		if !apply {
			return nil
		}
	}

	// run the apply command
	for _, c := range cmd.Root().Commands() {
		// set the apply command
		if c.Use == "apply" {
			// run the apply command
			if err := c.RunE(c, args); err != nil {
				return err
			}

			// call the post run command to cleanup
			return c.PostRunE(c, args)
		}
	}

	return nil
}

// VerifyInit is used to verify the init command has been run by checking if a config file exists.
func VerifyInit(cmd *cobra.Command, args []string, home string, output terminal.Outputer) error {
	// verify the config exists
	_, err := config.Load(home)
	if errors.Is(err, config.ErrNoConfigFile) {
		output.Info("Warning:", err.Error())

		// ask if the init command should run
		init, err := output.Confirm("Run `nitro init` now to create the config?", true, "")
		if err != nil {
			return err
		}

		// if init is false return nil
		if !init {
			return fmt.Errorf("You must run `nitro init` in order to add a site.")
		}

		// run the init command
		for _, c := range cmd.Root().Commands() {
			// set the init command
			if c.Use == "init" {
				if err := c.RunE(c, args); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
