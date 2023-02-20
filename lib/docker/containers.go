package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

func GetImages(ref string) ([]types.ImageSummary, error) {
	ref += "*"
	iFilter := filters.NewArgs(filters.Arg("reference", ref))
	return dockerClient.ImageList(
		context.Background(),
		types.ImageListOptions{
			Filters: iFilter,
		},
	)
}

func GetContainers(ref string) ([]types.Container, error) {
	cFilter := filters.NewArgs(filters.Arg("name", ref))
	return dockerClient.ContainerList(
		context.Background(),
		types.ContainerListOptions{
			All: true,
			Filters: cFilter,
		},
	)
}

func StartContainer(ref, name string) error {
	// just try to pull, if already present this will not be noticed
	err := MaybePullImage(ref)
	if err != nil {
		return err
	}

	// maybe stop here, by ignoring error
	// we do this to cleanup and stopped / done images
	// (typical of docker when a user restarts their computer)
	StopContainer(name)

	// now we can safely (re)create our container
	ret, err := dockerClient.ContainerCreate(
		context.Background(),

		// config
		// todo, maybe walk back versions? (dirty -> latest release)
		&container.Config{
			Image: ref,
		},

		// hostConfig
		&container.HostConfig{
			PublishAllPorts: true,
		},

		// netConfig
		nil,

		// todo, need to consider mac arm here?
		// platform
		nil,

		// name
		name,
	)

	if err != nil {
		return err
	}

	err = dockerClient.ContainerStart(
		context.Background(),
		ret.ID,
		types.ContainerStartOptions{},
	)

	return err
}

func StopContainer(name string) error {
	return dockerClient.ContainerRemove(
		context.Background(),
		name,
		types.ContainerRemoveOptions{ Force: true },
	)
}

func PullImage(ref string) error {
	opts := types.ImagePullOptions{}
	r, err := dockerClient.ImagePull( context.Background(), ref, opts)
	if err != nil {
		return err
	}

	_, err = dockerClient.ImageLoad(context.Background(), r, false)
	if err != nil {
		return err
	}

	return nil
}

func MaybePullImage(ref string) error {
	iFilter := filters.NewArgs(filters.Arg("reference", ref))
	images, err := dockerClient.ImageList(
		context.Background(),
		types.ImageListOptions{
			Filters: iFilter,
		},
	)
	if err != nil {
		return err
	}

	found := false
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == ref {
				found = true
				break
			}
		}
	}

	if !found {
		return PullImage(ref)
	}

	return nil
}

