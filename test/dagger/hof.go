package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

// so we don't have to pass these around everywhere
type runtime struct {
	ctx    context.Context
	client *dagger.Client
}

type stage func (*dagger.Container) (*dagger.Container, error)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	ctx := context.Background()

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	R := &runtime{
		ctx:    ctx,
		client: client,
	}

	var c *dagger.Container

	c, err = R.buildBase(nil)
	checkErr(err)

	c, err = R.loadCodeAndDeps(c)
	checkErr(err)

	c, err = R.buildHof(c)
	checkErr(err)

	c, d, err := R.buildHofMatrix(c)
	checkErr(err)

	ok, err := d.Export(R.ctx, ".")
	checkErr(err)
	if !ok {
		panic("unable to write matrix build outputs")
	}

	err = R.sanityTest(c)
	checkErr(err)

	err = R.dockerTest(c)
	checkErr(err)

}

func (R *runtime) loadCodeAndDeps(c *dagger.Container) (*dagger.Container, error) {
	// c = c.Pipeline("load")
	// setup mod cache
	modCache := R.client.CacheVolume("gomod")
	c = c.WithMountedCache("/go/pkg/mod", modCache)

	// load hof's code
	code := R.client.Host().Directory(".", dagger.HostDirectoryOpts{
		Exclude: []string{"cue.mod/pkg", "docs", "next"},
	})

	// get mods
	c = c.WithDirectory("/work", code, dagger.ContainerWithDirectoryOpts{
		Include: []string{"go.mod", "go.sums"},
	})
	c = c.WithExec([]string{"go", "mod", "tidy"})

	// build hof
	c = c.WithDirectory("/work", code)

	c = c.WithEnvVariable("CGO_ENABLED", "0")

	return c, nil
}

func (R *runtime) buildHof(c *dagger.Container) (*dagger.Container, error) {
	// c = c.Pipeline("build")
	c = c.WithExec([]string{"go", "build", "./cmd/hof"})
	// for testing
	c = c.WithExec([]string{"cp", "hof", "/usr/bin/hof"})

	return c, nil
}

func (R *runtime) buildHofMatrix(c *dagger.Container) (*dagger.Container, *dagger.Directory, error) {
	// the matrix
	geese := []string{"linux", "darwin"}
	goarches := []string{"amd64", "arm64"}

	// c = c.Pipeline("matrix")

	// build matrix for writing to host
	outputs := R.client.Directory()
	for _, goos := range geese {
		for _, goarch := range goarches {
			// create a directory for each OS and architecture
			path := fmt.Sprintf("build/%s/%s/", goos, goarch)

			// env vars
			build := c.
				WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch)

			// name the build
			// build = build.Pipeline(path)
			build = build.WithExec([]string{"go", "build", "-o", path, "./cmd/hof"})

			// add build to outputs
			outputs = outputs.WithDirectory(path, build.Directory(path))
		}
	}

	return c, outputs, nil
}

func (R *runtime) buildBase(c *dagger.Container) (*dagger.Container, error) {

	if c == nil {
		c = R.client.Container().From("golang:1.20")
	}

	// c = c.Pipeline("base")

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
	out, err := c.Stdout(R.ctx)
	if err != nil {
		fmt.Println(out)
		return c, err
	}

	// add docker CLI
	dockerCLI := R.client.Container().From("docker:24").
		File("/usr/local/bin/docker")

	c = c.WithFile("/usr/local/bin/docker", dockerCLI)

	// setup workdir
	c = c.WithWorkdir("/work")

	return c, nil
}

func (R *runtime) sanityTest(c *dagger.Container) error {
	t := c.WithExec([]string{"hof", "version"})
	// t = t.Pipeline("test/sanity")

	out, err := t.Stdout(R.ctx)
	if err != nil {
		fmt.Println(out)
		return err
	}

	return nil
}

func (R *runtime) dockerTest(c *dagger.Container) error {
	// c = c.Pipeline("test/flow")
	// add socket
	sock := R.client.Host().UnixSocket("/var/run/docker.sock")
	c = c.WithUnixSocket("/var/run/docker.sock", sock)

	c = c.WithExec([]string{"hof", "fmt", "info"})
	// c = c.WithExec([]string{"hof", "fmt", "pull", "v0.6.8-rc.5"})

	out, err := c.Stdout(R.ctx)
	if err != nil {
		fmt.Println(out)
		return err
	}

	fmt.Println(out)
	return nil
}
