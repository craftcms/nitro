package appcontainer

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/craftcms/nitro/pkg/bindmounts"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/dockerbuild"
	"github.com/craftcms/nitro/pkg/dockervolume"
	"github.com/craftcms/nitro/pkg/envvars"
	"github.com/craftcms/nitro/pkg/match"
	"github.com/craftcms/nitro/pkg/nginx"
	"github.com/craftcms/nitro/pkg/paths"
	"github.com/craftcms/nitro/pkg/proxycontainer"
	"github.com/craftcms/nitro/pkg/wsl"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
)

var Image = "docker.io/craftcms/nitro:%s"

type command struct {
	Name     string
	Commands []string
}

func StartOrCreate(home string, ctx context.Context, docker client.CommonAPIClient, cfg *config.Config, app config.App, networkID string) (string, error) {
	// check if nitro development is defined and override the image
	if _, ok := os.LookupEnv("NITRO_DEVELOPMENT"); ok {
		Image = "craftcms/nitro:%s"
	}

	// create a filter
	filter := filters.NewArgs()
	filter.Add("label", containerlabels.Host+"="+app.Hostname)

	// does the app use a custom dockerfile?
	if app.Dockerfile {
		// set the image name
		Image = fmt.Sprintf("%s:local", app.Hostname)

		imageFilter := filters.NewArgs()
		imageFilter.Add("reference", Image)

		// check if the image exists
		list, err := docker.ImageList(ctx, types.ImageListOptions{
			All:     true,
			Filters: imageFilter,
		})
		if err != nil {
			return "", err
		}

		if len(list) == 0 {
			path, err := paths.Clean(home, app.Path)
			if err != nil {
				return "", err
			}

			if err := dockerbuild.Build(os.Stdin, os.Stdout, path, Image); err != nil {
				return "", err
			}
		}
	}

	// look for the container for the app
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return "", err
	}

	// if there are no containers, create one
	if len(containers) == 0 {
		return create(ctx, docker, cfg, app, networkID)
	}

	// there is a container so get the first one
	c := containers[0]
	if c.State != "running" {
		if err := docker.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
			return "", err
		}
	}

	// get the containers details that include environment variables
	details, err := docker.ContainerInspect(ctx, c.ID)
	if err != nil {
		return "", err
	}

	// if the container is out of date
	if !match.App(home, app, details, cfg.Blackfire) {
		fmt.Print("- updating… ")

		// stop container
		if err := docker.ContainerStop(ctx, c.ID, nil); err != nil {
			return "", err
		}

		// remove container
		if err := docker.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{}); err != nil {
			return "", err
		}

		return create(ctx, docker, cfg, app, networkID)
	}

	return c.ID, nil
}

func ContainerPath(app config.App) string {
	// trim trailing slashes
	webroot := strings.TrimRight(app.Webroot, "/")

	// is there a path separator?
	if strings.Contains(webroot, "/") {
		parts := strings.Split(webroot, "/")

		if len(parts) >= 2 {
			return strings.Join(parts[:len(parts)-1], "/")
		}
	}

	return "/app"
}

func create(ctx context.Context, docker client.CommonAPIClient, cfg *config.Config, app config.App, networkID string) (string, error) {
	// create the container
	image := fmt.Sprintf(Image, app.PHPVersion)
	// if the app uses a dockerfile, make the image based on the hostname
	if app.Dockerfile {
		image = fmt.Sprintf("%s:local", app.Hostname)
	}

	// pull the image if we are not in a development environment
	_, dev := os.LookupEnv("NITRO_DEVELOPMENT")
	if !dev {
		rdr, err := docker.ImagePull(ctx, image, types.ImagePullOptions{All: false})
		if err != nil {
			return "", fmt.Errorf("unable to pull the image, %w", err)
		}

		buf := &bytes.Buffer{}
		if _, err := buf.ReadFrom(rdr); err != nil {
			return "", fmt.Errorf("unable to read output from pulling image %s, %w", image, err)
		}
	}

	// add the app itself and any aliases to the extra hosts
	extraHosts := []string{fmt.Sprintf("%s:%s", app.Hostname, "127.0.0.1")}
	for _, s := range app.Aliases {
		extraHosts = append(extraHosts, fmt.Sprintf("%s:%s", s, "127.0.0.1"))
	}

	// check if this is linux specific
	if runtime.GOOS == "linux" && !wsl.IsWSL() {
		extraHosts = append(extraHosts, fmt.Sprintf("%s:%s", "host.docker.internal", "host-gateway"))
	}

	// get the app environment variables
	envs, err := envvars.ForApp(app, "host.docker.internal")
	if err != nil {
		return "", err
	}

	// does the config have blackfire credentials
	if app.Blackfire {
		// grab the credentials from the config
		credentials, err := envvars.GetBlackfireClientCredentials(*cfg)
		if err != nil {
			return "", err
		}

		// add the client credentials
		envs = append(envs, credentials...)

		// set the agent socket to use the service container
		envs = append(envs, "BLACKFIRE_AGENT_SOCKET=tcp://blackfire.service.nitro:8307")
	}

	// create labels for the volume
	volumeLabels := containerlabels.ForAppVolume(app)

	// look for an existing volume with the app hostname, otherwise create it
	if err := dockervolume.CreateIfEmpty(ctx, docker, app.Hostname, volumeLabels); err != nil {
		return "", err
	}

	// look for an existing volume with the app hostname + nginx, otherwise create it
	if err := dockervolume.CreateIfEmpty(ctx, docker, fmt.Sprintf("%s-nginx", app.Hostname), nil); err != nil {
		return "", err
	}

	// determine if the site has any excludes
	binds, err := bindmounts.ForApp(app, cfg.HomeDirectory)
	if err != nil {
		return "", err
	}

	// set the labels
	labels := containerlabels.ForApp(app)
	// create the container
	resp, err := docker.ContainerCreate(
		ctx,
		&container.Config{
			Image:    image,
			Labels:   labels,
			Env:      envs,
			Hostname: app.Hostname,
		},
		&container.HostConfig{
			Binds:      binds,
			ExtraHosts: extraHosts,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: proxycontainer.VolumeName,
					Target: proxycontainer.VolumeTarget,
				},
				{
					Type:   mount.TypeVolume,
					Source: app.Hostname,
					Target: "/home/nitro",
				},
				{
					Type:   mount.TypeVolume,
					Source: fmt.Sprintf("%s-nginx", app.Hostname),
					Target: "/etc/nginx/sites-available/",
				},
			},
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"nitro-network": {
					NetworkID: networkID,
				},
			},
		},
		nil,
		app.Hostname,
	)
	if err != nil {
		return "", fmt.Errorf("unable to create the container, %w", err)
	}

	// start the container
	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("unable to start the container, %w", err)
	}

	restart := false

	// post installation commands
	var commands []command
	// check for a custom root and copy the template to the container
	if app.Webroot != "web" {
		// we need to restart the container to take effect
		restart = true

		// create the nginx file
		conf := nginx.Generate(app.Webroot)

		// create the temp file
		tr, err := archive.Generate(app.Hostname, conf)
		if err != nil {
			return "", err
		}

		// copy the file into the container
		if err := docker.CopyToContainer(ctx, resp.ID, "/tmp", tr, types.CopyToContainerOptions{AllowOverwriteDirWithFile: false}); err != nil {
			return "", err
		}

		commands = append(commands, command{Commands: []string{"cp", fmt.Sprintf("/tmp/%s", app.Hostname), "/etc/nginx/sites-available/default"}})
		commands = append(commands, command{Commands: []string{"chmod", "0644", "/etc/nginx/sites-available/default"}})
	}

	// check if there are custom extensions, NOTE: extensions require a container restart
	for _, ext := range app.Extensions {
		// we need to restart the container
		restart = true

		commands = append(commands, command{
			Name:     "installing-" + ext + "-extension",
			Commands: []string{"apt", "install", "--yes", "–no-install-recommends", fmt.Sprintf("php%s-%s", app.PHPVersion, ext)},
		})
	}

	// run the commands
	for _, c := range commands {
		// create the exec
		exec, err := docker.ContainerExecCreate(ctx, resp.ID, types.ExecConfig{
			User:         "root",
			AttachStdout: true,
			AttachStderr: true,
			Tty:          false,
			Cmd:          c.Commands,
		})
		if err != nil {
			return "", err
		}

		// attach to the container
		attach, err := docker.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{
			Tty: false,
		})
		if err != nil {
			return "", err
		}
		defer attach.Close()

		// if the option is for a php extension, don't show output
		if strings.Contains(c.Name, "-extension") {
			// read the output to pull the image
			fmt.Print("installing ", c.Commands[len(c.Commands)-1], "… ")

			buf := &bytes.Buffer{}
			if _, err := buf.ReadFrom(attach.Reader); err != nil {
				return "", fmt.Errorf("unable to read output from container exec attach, %w", err)
			}
		} else {
			// show the output to stdout and stderr
			if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, attach.Reader); err != nil {
				return "", fmt.Errorf("unable to copy the output of container, %w", err)
			}
		}

		// start the exec
		if err := docker.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{}); err != nil {
			return "", fmt.Errorf("unable to start the container, %w", err)
		}

		// wait for the container exec to complete
		waiting := true
		for waiting {
			resp, err := docker.ContainerExecInspect(ctx, exec.ID)
			if err != nil {
				return "", err
			}

			waiting = resp.Running
		}

		// start the container
		if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return "", fmt.Errorf("unable to start the container, %w", err)
		}
	}

	// restart the container if there is a custom extension
	if restart {
		if err := docker.ContainerRestart(ctx, resp.ID, nil); err != nil {
			return "", err
		}
	}

	return resp.ID, nil
}
