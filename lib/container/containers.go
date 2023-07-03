package container

import (
	"context"
	"time"
)

func GetImages(ref string) ([]Image, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return rt.Images(ctx, Ref(ref))
}

func GetContainers(name string) ([]Container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return rt.Containers(ctx, Name(name))
}

func StartContainer(ref, name string, env []string, replace bool) error {
	if replace {
		StopContainer(name)
	}

	return rt.Run(context.Background(), Ref(ref), Params{
		Name: Name(name),
		Env:  env,
	})
}

func StopContainer(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return rt.Remove(ctx, Name(name))
}

func PullImage(ref string) error {
	return rt.Pull(context.Background(), Ref(ref))
}
