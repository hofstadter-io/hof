package main

import (
	"context"
	"os"

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

	source := client.Host().Directory(".", dagger.HostDirectoryOpts{
		Exclude: []string{"cue.mod/pkg", "docs", "next"},
	})

	R := &Runtime{
		Ctx:    ctx,
		Client: client,
	}

	// docker := R.DockerImage()
	golang := R.GolangImage()

	buildr := golang.Pipeline("builder")
	buildr = buildr.WithDirectory("/work", source)
	buildr = buildr.WithExec([]string{"go", "build", "./cmd/hof"})

	valid := buildr.Pipeline("validate")
	valid = valid.WithExec([]string{"./hof", "version"})


	final := valid
	final.Sync(ctx)
	// final.Stdout(ctx)
}

func (R *Runtime) DockerImage() (*dagger.Container) {
	d := R.Client.Container().From("docker:24")

	return d
}

func (R *Runtime) GolangImage() (*dagger.Container) {

	c := R.Client.Container().From("golang:1.20")

	// setup mod cache
	modCache := R.Client.CacheVolume("gomod")
	c = c.WithMountedCache("/go/pkg/mod", modCache)

	// setup build cache
	buildCache := R.Client.CacheVolume("go-build")
	c = c.WithMountedCache("/root/.cache/go-build", buildCache)

	// add tools
	c = R.AddDockerCLI(c)

	// setup workdir
	c = c.WithWorkdir("/work")

	return c
}

func (R *Runtime) AddDockerCLI(c *dagger.Container) (*dagger.Container) {
	dockerCLI := R.Client.Container().From(dockerVer).
		File("/usr/local/bin/docker")

	c = c.WithFile("/usr/local/bin/docker", dockerCLI)

	return c
}
