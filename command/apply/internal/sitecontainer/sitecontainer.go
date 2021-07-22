package sitecontainer

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/craftcms/nitro/command/apply/internal/match"
	"github.com/craftcms/nitro/command/apply/internal/nginx"
	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/craftcms/nitro/pkg/proxycontainer"
	"github.com/craftcms/nitro/pkg/wsl"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/stdcopy"
)

type command struct {
	Name     string
	Commands []string
}

var (
	// Image is the image used for sites, with the PHP version
	Image = "craftcms/nitro:%s"
)

// StartOrCreate is responsible for finding a sites existing container or creating a new one based on the values from the configuration file.
func StartOrCreate(ctx context.Context, docker client.CommonAPIClient, home, networkID string, site config.Site, cfg *config.Config) (string, error) {
	// set filters for the container
	filter := filters.NewArgs()
	filter.Add("label", containerlabels.Host+"="+site.Hostname)

	// look for a container for the site
	containers, err := docker.ContainerList(ctx, types.ContainerListOptions{All: true, Filters: filter})
	if err != nil {
		return "", fmt.Errorf("error getting a list of containers")
	}

	// if there are no containers we need to create one
	if len(containers) == 0 {
		return create(ctx, docker, home, networkID, site, cfg)
	}

	// there is a container, so inspect it and make sure it matched
	container := containers[0]

	if container.State != "running" {
		if err := docker.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
			return "", err
		}
	}

	// get the containers details that include environment variables
	details, err := docker.ContainerInspect(ctx, container.ID)
	if err != nil {
		return "", err
	}

	// if the container is out of date
	if !match.Site(home, site, details, cfg.Blackfire) {
		fmt.Print("- updating… ")

		// stop container
		if err := docker.ContainerStop(ctx, container.ID, nil); err != nil {
			return "", err
		}

		// remove container
		if err := docker.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
			return "", err
		}

		return create(ctx, docker, home, networkID, site, cfg)
	}

	return container.ID, nil
}

func create(ctx context.Context, docker client.CommonAPIClient, home, networkID string, site config.Site, cfg *config.Config) (string, error) {
	// create the container
	image := fmt.Sprintf(Image, site.Version)

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

	// get the sites path
	path, err := site.GetAbsPath(home)
	if err != nil {
		return "", err
	}

	// add the site itself and any aliases to the extra hosts
	extraHosts := []string{fmt.Sprintf("%s:%s", site.Hostname, "127.0.0.1")}
	for _, s := range site.Aliases {
		extraHosts = append(extraHosts, fmt.Sprintf("%s:%s", s, "127.0.0.1"))
	}

	// check if this is linux specific
	if runtime.GOOS == "linux" && !wsl.IsWSL() {
		extraHosts = append(extraHosts, fmt.Sprintf("%s:%s", "host.docker.internal", "host-gateway"))
	}

	// get the sites environment variables
	envs := site.AsEnvs("host.docker.internal")

	// does the config have blackfire credentials
	if cfg.Blackfire.ServerID != "" {
		envs = append(envs, "BLACKFIRE_SERVER_ID="+cfg.Blackfire.ServerID)
	}

	if cfg.Blackfire.ServerToken != "" {
		envs = append(envs, "BLACKFIRE_SERVER_TOKEN="+cfg.Blackfire.ServerToken)
	}

	// look for an existing volume with the sites hostname, otherwise create it
	filter := filters.NewArgs()
	filter.Add("name", site.Hostname)
	volumeResp, err := docker.VolumeList(ctx, filter)
	if err != nil {
		return "", err
	}

	if len(volumeResp.Volumes) == 0 {
		labels := containerlabels.ForSiteVolume(site)

		if _, err := docker.VolumeCreate(ctx, volume.VolumeCreateBody{Driver: "local", Name: site.Hostname, Labels: labels}); err != nil {
			return "", err
		}
	}

	// look for an existing volume with the sites hostname + nginx, otherwise create it
	nginxFilter := filters.NewArgs()
	nginxFilter.Add("name", fmt.Sprintf("%s-nginx", site.Hostname))
	nginxVol, err := docker.VolumeList(ctx, nginxFilter)
	if err != nil {
		return "", err
	}

	if len(nginxVol.Volumes) == 0 {
		if _, err := docker.VolumeCreate(ctx, volume.VolumeCreateBody{Driver: "local", Name: fmt.Sprintf("%s-nginx", site.Hostname)}); err != nil {
			return "", err
		}
	}

	// set the labels
	labels := containerlabels.ForSite(site)
	// create the container
	resp, err := docker.ContainerCreate(
		ctx,
		&container.Config{
			Image:    image,
			Labels:   labels,
			Env:      envs,
			Hostname: site.Hostname,
		},
		&container.HostConfig{
			Binds:      []string{fmt.Sprintf("%s:/app:rw", path)},
			ExtraHosts: extraHosts,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeVolume,
					Source: proxycontainer.VolumeName,
					Target: proxycontainer.VolumeTarget,
				},
				{
					Type:   mount.TypeVolume,
					Source: site.Hostname,
					Target: "/home/nitro",
				},
				{
					Type:   mount.TypeVolume,
					Source: fmt.Sprintf("%s-nginx", site.Hostname),
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
		site.Hostname,
	)
	if err != nil {
		return "", fmt.Errorf("unable to create the container, %w", err)
	}

	// start the container
	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("unable to start the container, %w", err)
	}

	// post installation commands
	var commands []command

	// check for a custom root and copt the template to the container
	if site.Webroot != "web" {
		// create the nginx file
		conf := nginx.Generate(site.Webroot)

		// create the temp file
		tr, err := archive.Generate("default.conf", conf)
		if err != nil {
			return "", err
		}

		// copy the file into the container
		if err := docker.CopyToContainer(ctx, resp.ID, "/tmp", tr, types.CopyToContainerOptions{AllowOverwriteDirWithFile: false}); err != nil {
			return "", err
		}

		commands = append(commands, command{Commands: []string{"cp", "/tmp/default.conf", "/etc/nginx/conf.d/default.conf"}})
		commands = append(commands, command{Commands: []string{"chmod", "0644", "/etc/nginx/conf.d/default.conf"}})
	}

	// check if there are custom extensions
	for _, ext := range site.Extensions {
		c := command{
			Name:     "installing-" + ext + "-extension",
			Commands: []string{"apt", "install", "--yes", "–no-install-recommends", fmt.Sprintf("php%s-%s", site.Version, ext)},
		}
		commands = append(commands, c)
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

	return resp.ID, nil
}
