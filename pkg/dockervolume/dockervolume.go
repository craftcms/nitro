package dockervolume

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

func CreateIfEmpty(ctx context.Context, docker client.VolumeAPIClient, name string, labels map[string]string) error {
	f := filters.NewArgs()
	f.Add("name", name)
	vol, err := docker.VolumeList(ctx, f)
	if err != nil {
		return err
	}

	// if there are no volumes, create one
	if len(vol.Volumes) == 0 {
		if _, err := docker.VolumeCreate(ctx, volume.VolumeCreateBody{
			Driver: "local",
			Name:   name,
			Labels: labels,
		}); err != nil {
			return err
		}
	}

	return nil
}
