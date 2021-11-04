package actions

import (
	"fmt"

	"github.com/craftcms/nitro/pkg/config"
	"github.com/craftcms/nitro/pkg/containerlabels"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

// AppToContainerConfig
func AppToContainerConfig(app config.App) *container.Config {
	// check for custom dockerfile
	image := fmt.Sprintf("craftcms/nitro:%s", app.PHPVersion)
	if app.Dockerfile {
		image = fmt.Sprintf("%s:local", app.GetHostname())
	}

	// get the hostname
	hostname := app.GetHostname()

	// create the labels
	labels := map[string]string{
		containerlabels.Nitro:   "true",
		containerlabels.Host:    hostname,
		containerlabels.Webroot: app.Webroot,
		containerlabels.Type:    "app",
	}

	// get the environment variables for the app
	env := app.AsEnvs("host.docker.internal")

	return &container.Config{
		Image:    image,
		Labels:   labels,
		Env:      env,
		Hostname: hostname,
	}
}

func AppToContainerHostConfig(app config.App) *container.HostConfig {
	return &container.HostConfig{
		Binds:      nil,
		ExtraHosts: nil,
		Mounts:     nil,
	}
}

func AppToContainerNetworkingConfig(app config.App, id string) *network.NetworkingConfig {
	return &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"nitro-network": {
				NetworkID: id,
			},
		},
	}
}
