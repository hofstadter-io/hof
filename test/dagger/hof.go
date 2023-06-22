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

func (R *Runtime) LocalCodeAndDeps(c *dagger.Container) (*dagger.Container, error) {
	c = c.Pipeline("load")
	// setup mod cache
	modCache := R.Client.CacheVolume("gomod")
	c = c.WithMountedCache("/go/pkg/mod", modCache)

	// load hof's code
	code := R.Client.Host().Directory(".", dagger.HostDirectoryOpts{
		Exclude: []string{"cue.mod/pkg", "docs", "next"},
	})

	// get mods
	c = c.WithDirectory("/work", code, dagger.ContainerWithDirectoryOpts{
		Include: []string{"go.mod", "go.sums"},
	})
	c = c.WithExec([]string{"go", "mod", "tidy"})

	// add full code
	c = c.WithDirectory("/work", code)

	// ensure cgo is disabled everywhere
	c = c.WithEnvVariable("CGO_ENABLED", "0")

	return c, nil
}

func (R *Runtime) BuildHof(c *dagger.Container) (*dagger.Container, error) {
	c = c.Pipeline("build")
	c = c.WithExec([]string{"go", "build", "./cmd/hof"})
	// for testing
	c = c.WithExec([]string{"cp", "hof", "/usr/bin/hof"})

	return c, nil
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

func (R *Runtime) BuildBase(c *dagger.Container) (*dagger.Container, error) {

	if c == nil {
		c = R.Client.Container().From("golang:1.20")
	}

	c = c.Pipeline("base")

	// install packages
	c = c.WithExec([]string{
		"bash", "-c",
		`
		apt-get update -y && \
		apt-get install -y \
		tree && \
		apt search docker
		`,
	})
	out, err := c.Stdout(R.Ctx)
	if err != nil {
		fmt.Println(out)
		return c, err
	}

	// add docker CLI
	dockerCLI := R.Client.Container().From("docker:24").
		File("/usr/local/bin/docker")

	c = c.WithFile("/usr/local/bin/docker", dockerCLI)

	// setup workdir
	c = c.WithWorkdir("/work")

	return c, nil
}

func (R *Runtime) SanityTest(c *dagger.Container) error {
	t := c.Pipeline("test/sanity")
	t = t.WithExec([]string{"hof", "version"})

	out, err := t.Stdout(R.Ctx)
	if err != nil {
		fmt.Println(out)
		return err
	}

	return nil
}

func (R *Runtime) DockerTest(c *dagger.Container) error {
	t := c.Pipeline("test/flow")
	// add socket
	sock := R.Client.Host().UnixSocket("/var/run/docker.sock")
	t = t.WithUnixSocket("/var/run/docker.sock", sock)

	t = t.WithExec([]string{"hof", "fmt", "info"})
	// c = c.WithExec([]string{"hof", "fmt", "pull", "v0.6.8-rc.5"})

	out, err := c.Stdout(R.Ctx)
	if err != nil {
		fmt.Println(out)
		return err
	}

	fmt.Println(out)
	return nil
}
