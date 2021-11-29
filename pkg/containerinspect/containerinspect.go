package containerinspect

import (
	"context"
	"sort"

	"github.com/docker/docker/client"
)

type Inspecter interface {
	Inspect(ctx context.Context, docker client.ContainerAPIClient, container string) (Info, error)
}

type Info struct {
	User             string
	WorkingDirectory string
	Ports            []int
}

func Inspect(ctx context.Context, docker client.ContainerAPIClient, container string) (Info, error) {
	resp, err := docker.ContainerInspect(ctx, container)
	if err != nil {
		return Info{}, err
	}

	info := Info{}

	info.User = resp.Config.User
	info.WorkingDirectory = resp.Config.WorkingDir

	for p := range resp.Config.ExposedPorts {
		info.Ports = append(info.Ports, p.Int())
	}

	sort.Slice(info.Ports, func(p, q int) bool {
		return info.Ports[p] < info.Ports[q]
	})

	return info, nil
}
