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


func (R *Runtime) Hack(c *dagger.Container) (*dagger.Container, error) {
	t := c.Pipeline("hack")

	// dev time function, do whatever here

	return t, nil
}

func (R *Runtime) BaseContainer() (*dagger.Container) {

	c := R.Client.Container().From("golang:1.20")

	c = c.Pipeline("base")

	// setup workdir
	c = c.WithWorkdir("/work")

	// add tools
	c = R.AddDockerCLI(c)

	return c
}

func (R *Runtime) RuntimeContainer(builder *dagger.Container) (*dagger.Container) {
	hof := builder.File("/work/hof")

	c := R.BaseContainer()
	c = c.Pipeline("hof/runtime")
	c = c.WithFile("/usr/local/bin/hof", hof)
	
	return c
}

func (R *Runtime) WithCodeAndDeps(c *dagger.Container, source *dagger.Directory) (*dagger.Container) {
	c = c.Pipeline("hof/load")

	// setup mod cache
	modCache := R.Client.CacheVolume("gomod")
	c = c.WithMountedCache("/go/pkg/mod", modCache)

	// get mods
	c = c.WithDirectory("/work", source, dagger.ContainerWithDirectoryOpts{
		Include: []string{"go.mod", "go.sums"},
	})
	c = c.WithExec([]string{"go", "mod", "download"})

	// add full code
	c = c.WithDirectory("/work", source)

	return c
}

func (R *Runtime) BuildHof(c *dagger.Container) (*dagger.Container) {
	c = c.Pipeline("hof/build")
	c = c.WithEnvVariable("CGO_ENABLED", "0")
	c = c.WithExec([]string{"go", "build", "./cmd/hof"})
	c = c.WithExec([]string{"cp", "hof", "/usr/local/bin/hof"})
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
