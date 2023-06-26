package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
)

const dockerVer = "docker:24"

type Runtime struct {
	Ctx    context.Context
	Client *dagger.Client
}

func main() {
	ctx := context.Background()

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	R := &Runtime{
		Ctx:    ctx,
		Client: client,
	}

	// docker := R.DockerImage()
	gcloud := R.GcloudImage()

	t := gcloud
	t = t.WithExec([]string{"gcloud", "version"})
	t = t.WithExec([]string{"gcloud", "config", "list"})


	final := t
	final.Sync(ctx)
	// final.Stdout(ctx)
}

func (R *Runtime) GcloudImage() (*dagger.Container) {

	cfg, err := os.UserConfigDir()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	d := R.Client.Host().Directory(filepath.Join(cfg, "gcloud"))

	c := R.Client.Container().From("google/cloud-sdk")
	c = c.WithEnvVariable("CLOUDSDK_CONFIG", "/gcloud/config")
	c = c.WithDirectory("/gcloud/config", d)

	return c
}
