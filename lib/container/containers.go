package container

import (
	"context"
	"time"
)

// these should all take a context object

func GetImages(ref string) ([]Image, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return rt.Images(ctx, Ref(ref))
}

func GetContainers(name string) ([]Container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return rt.Containers(ctx, Name(name))
}

func StartContainer(ref string, params *Params) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return rt.Run(ctx, Ref(ref), params)
}

func StopContainer(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return rt.Remove(ctx, Name(name))
}

func PullImage(ref string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	return rt.Pull(ctx, Ref(ref))
}

func LoadTarball(ctx context.Context, content []byte) error {
	return rt.Load(ctx, "", []byte(content))
}
