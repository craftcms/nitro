package add

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/labels"
	"github.com/craftcms/nitro/pkg/phpversions"
	"github.com/craftcms/nitro/pkg/terminal"
	"github.com/craftcms/nitro/pkg/validate"
	"github.com/craftcms/nitro/pkg/webroot"
)

// TODO - prompt user for the database engine and new database
// TODO - edit the env file for the user

const exampleText = `  # add the current project as a site
  nitro add

  # add a directory as the site
  nitro add my-project`

func NewCommand(home string, docker client.CommonAPIClient, output terminal.Outputer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Short:   "Add a site",
		Example: exampleText,
		PostRunE: func(cmd *cobra.Command, args []string) error {
			// ask if the apply command should run
			apply, err := output.Confirm("Apply changes now", true, "?")
			if err != nil {
				return err
			}

			// if apply is false return nil
			if !apply {
				return nil
			}

			// run the apply command
			for _, c := range cmd.Parent().Commands() {
				// set the apply command
				if c.Use == "apply" {
					return c.RunE(c, args)
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			output.Info("Adding site…")

			// get the current working directory
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			// get working directory or provided arg
			var dir string
			switch len(args) {
			case 1:
				dir = filepath.Join(wd, args[0])
			default:
				dir = filepath.Clean(wd)
			}

			// create a new site
			site := config.Site{}

			// get the hostname from the directory
			sp := strings.Split(dir, string(os.PathSeparator))
			site.Hostname = sp[len(sp)-1]

			// append the default TLD if there are no periods in the path name
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
				return err
			}

			// set the input as the hostname
			site.Hostname = hostname

			output.Success("setting hostname to", site.Hostname)

			// set the sites directory but make the path relative
			site.Path = strings.Replace(dir, home, "~", 1)

			output.Success("adding site", site.Path)

			// get the web directory
			found, err := webroot.Find(dir)
			if err != nil {
				return err
			}

			if found == "" {
				found = "web"
			}

			// set the webroot
			site.Webroot = found

			// prompt for the webroot
			root, err := output.Ask("Enter the webroot for the site", site.Webroot, ":", nil)
			if err != nil {
				return err
			}

			site.Webroot = root

			output.Success("using webroot", site.Webroot)

			// prompt for the php version
			versions := phpversions.Versions
			selected, err := output.Select(cmd.InOrStdin(), "Choose a PHP version: ", versions)
			if err != nil {
				return err
			}

			// set the version of php
			site.Version = versions[selected]

			output.Success("setting PHP version", site.Version)

			// load the config
			cfg, err := config.Load(home)
			if err != nil {
				return err
			}

			// add the site to the config
			if err := cfg.AddSite(site); err != nil {
				return err
			}

			output.Pending("saving file")

			// save the config file
			if err := cfg.Save(); err != nil {
				output.Warning()

				return err
			}

			output.Done()

			output.Info("Site added 🌍")

			// prompt for a database
			if err := promptForDatabase(cmd.Context(), docker, output); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func promptForDatabase(ctx context.Context, docker client.CommonAPIClient, output terminal.Outputer) error {
	fmt.Print("Add a database for the site? [Y/n]? ")

	s := bufio.NewScanner(os.Stdin)
	s.Split(bufio.ScanLines)

	var confirm bool
	for s.Scan() {
		txt := strings.TrimSpace(s.Text())

		switch txt {
		// if its no
		case "n", "N", "no", "No", "NO":
			confirm = false
		default:
			confirm = true
		}

		break
	}
	if err := s.Err(); err != nil {
		return err
	}

	if !confirm {
		return nil
	}

	// add filters to show only the environment and database containers
	filter := filters.NewArgs()
	filter.Add("label", labels.Nitro)
	filter.Add("label", labels.Type+"=database")

	// get a list of all the databases
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{Filters: filter})
	if err != nil {
		return err
	}

	// sort containers by the name
	sort.SliceStable(containers, func(i, j int) bool {
		return containers[i].Names[0] < containers[j].Names[0]
	})

	// get all of the containers as a list
	var engineOpts []string
	for _, c := range containers {
		engineOpts = append(engineOpts, strings.TrimLeft(c.Names[0], "/"))
	}

	// prompt the user for the engine to add the database
	var containerID, databaseEngine string
	selected, err := output.Select(os.Stdin, "Select the database engine: ", engineOpts)
	if err != nil {
		return err
	}

	// set the container id and db engine
	containerID = containers[selected].ID
	databaseEngine = containers[selected].Labels[labels.DatabaseCompatibility]
	if containerID == "" {
		return fmt.Errorf("unable to get the container")
	}

	// ask the user for the database to create
	db, err := output.Ask("Enter the new database name", "", ":", nil)
	if err != nil {
		return err
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
		return err
	}

	// attach to the container
	resp, err := docker.ContainerExecAttach(ctx, e.ID, types.ExecStartCheck{
		Tty: false,
	})
	if err != nil {
		return err
	}
	defer resp.Close()

	// start the exec
	if err := docker.ContainerExecStart(ctx, e.ID, types.ExecStartCheck{}); err != nil {
		return fmt.Errorf("unable to start the container, %w", err)
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
			return err
		}

		// attach to the container
		resp, err := docker.ContainerExecAttach(ctx, e.ID, types.ExecStartCheck{
			Tty: false,
		})
		if err != nil {
			return err
		}
		defer resp.Close()

		// start the exec
		if err := docker.ContainerExecStart(ctx, e.ID, types.ExecStartCheck{}); err != nil {
			return fmt.Errorf("unable to start the container, %w", err)
		}

		// wait for the container exec to complete
		waiting := true
		for waiting {
			resp, err := docker.ContainerExecInspect(ctx, e.ID)
			if err != nil {
				return err
			}

			waiting = resp.Running
		}
	}

	output.Done()

	output.Info("Database added 💪")

	// TODO(jasonmccallister) edit the env

	return nil

}
