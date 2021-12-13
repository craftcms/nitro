package imageinspect

import (
	"context"
	"fmt"
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

// Inspect inspects an image and looks for details required for interacting with a container in nitro.
func Inspect(ctx context.Context, docker client.ImageAPIClient, image string) (*Info, error) {
	resp, _, err := docker.ImageInspectWithRaw(ctx, image)
	if err != nil {
		return nil, err
	}

	info := Info{
		User:             resp.Config.User,
		WorkingDirectory: resp.Config.WorkingDir,
	}

	for p := range resp.Config.ExposedPorts {
		info.Ports = append(info.Ports, p.Int())
	}

	sort.Slice(info.Ports, func(p, q int) bool {
		return info.Ports[p] < info.Ports[q]
	})

	// validate the info
	if info.User == "" {
		return nil, fmt.Errorf("missing the USER argument from the image")
	}

	// ensure there is a working directory
	if info.WorkingDirectory == "" {
		return nil, fmt.Errorf("missing the WORKDIR argument from the image")
	}

	// verify there are ports
	if len(info.Ports) == 0 {
		return nil, fmt.Errorf("missing the EXPOSE argument from the image")
	}

	return &info, nil
}
