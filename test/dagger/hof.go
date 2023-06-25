package dagger

import (
	"context"
	"fmt"
	"strings"

	"dagger.io/dagger"
)

// so we don't have to pass these around everywhere
type Runtime struct {
	Ctx    context.Context
	Client *dagger.Client
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

func (R *Runtime) RuntimeContainer(builder *dagger.Container) (*dagger.Container) {
	hof := builder.File("hof")

	c := R.GolangImage()
	c = c.WithFile("/usr/local/bin/hof", hof)
	c = c.Pipeline("hof/runtime")
	
	return c
}

func (R *Runtime) FetchDeps(c *dagger.Container, source *dagger.Directory) (*dagger.Container) {
	c = c.Pipeline("hof/deps")

	// get deps
	c = c.WithDirectory("/work", source, dagger.ContainerWithDirectoryOpts{
		Include: []string{"go.mod", "go.sums"},
	})
	c = c.WithExec([]string{"go", "mod", "download"})

	// c = c.WithDirectory("/work", source)
	return c
}

func (R *Runtime) BuildHof(c *dagger.Container, source *dagger.Directory) (*dagger.Container) {
	c = c.Pipeline("hof/build")

	// exclude files we don't need so we can avoid cache misses?
	c = c.WithDirectory("/work", source, dagger.ContainerWithDirectoryOpts{
		Exclude: []string{
			"changelogs",
			"ci",
			"docs",
			"hack",
			"images",
			"notes",
			"test", 
		},
	})

	c = c.WithEnvVariable("CGO_ENABLED", "0")

	c = c.WithExec([]string{"go", "build", "./cmd/hof"})
	return c
}

func (R *Runtime) BuildHofMatrix(c *dagger.Container) (*dagger.Directory, error) {
	// the matrix
	geese := []string{"linux", "darwin"}
	goarches := []string{"amd64", "arm64"}

	// c = c.Pipeline("matrix")

	outputs := R.Client.Directory()
	outputs = outputs.Pipeline("outputs")

	// build matrix for writing to host
	for _, goos := range geese {
		for _, goarch := range goarches {
			// create a directory for each OS and architecture
			path := fmt.Sprintf("build/%s/%s/", goos, goarch)

			// name the build
			build := c.Pipeline(strings.TrimSuffix(path, "/"))

			// set local env vars
			build = build.
				WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch)

			// run the build
			build = build.WithExec([]string{"go", "build", "-o", path, "./cmd/hof"})

			// add build to outputs
			outputs = outputs.WithDirectory(path, build.Directory(path))
		}
	}

	return outputs, nil
}

func (R *Runtime) HofVersion(c *dagger.Container) error {
	t := c.Pipeline("hof/version")
	t = t.WithExec([]string{"hof", "version"})

	_, err := t.Sync(R.Ctx)
	return err
}
