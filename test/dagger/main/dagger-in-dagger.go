package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"dagger.io/dagger"
)

// meta...?
var meta int64 = 0

func main() {
	fmt.Println(os.Args)

	if len(os.Args) == 2 {
		m := os.Args[1]
		n, err := strconv.ParseInt(m, 10, 64)
		if err != nil {
			panic(err)
		}
		meta = n
	}

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

	c, err := R.buildContainer()
	if err != nil {
		panic(err)
	}

	err = R.sanityTest(c)
	if err != nil {
		panic(err)
	}

	err = R.dockerTest(c)
	if err != nil {
		panic(err)
	}

	if meta > 0 {
		err = R.inception(c)
		if err != nil {
			panic(err)
		}
	}
}

type runtime struct {
	ctx    context.Context
	client *dagger.Client
}

type processor func(c *dagger.Container) (*dagger.Container, error)

func (R *runtime) buildContainer() (*dagger.Container, error) {

	c := R.client.Container().From("golang:1.20")

	var err error
	for _, fn := range []processor{
		R.addPackages,
		R.addDockerCLI,
	} {
		c, err = fn(c)
		if err != nil {
			return c, err
		}
	}

	return c, nil
}

func (R *runtime) addPackages(c *dagger.Container) (*dagger.Container, error) {
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
		return c, err
	}
	_, err = c.Sync(R.ctx)
	if err != nil {
		return c, err
	}
	fmt.Println(out)

	return c, nil
}

func (R *runtime) addDockerCLI(c *dagger.Container) (*dagger.Container, error) {
	dockerCLI := R.client.Container().From("docker:24").
		File("/usr/local/bin/docker")

	c = c.WithFile("/usr/local/bin/docker", dockerCLI)

	return c, nil
}

func (R *runtime) sanityTest(c *dagger.Container) error {
	t := c.WithExec([]string{"tree", "--version"})

	out, err := t.Stdout(R.ctx)
	if err != nil {
		fmt.Println(out)
		return err
	}

	return nil
}

func (R *runtime) dockerTest(c *dagger.Container) error {
	sock := R.client.Host().UnixSocket("/var/run/docker.sock")

	c = c.WithUnixSocket("/var/run/docker.sock", sock)

	c = c.WithExec([]string{"docker", "info"})

	out, err := c.Stdout(R.ctx)
	if err != nil {
		fmt.Println(out)
		return err
	}

	fmt.Println(out)
	return nil
}

func (R *runtime) inception(c *dagger.Container) error {
	if meta == 0 {
		fmt.Println("you found the bottom of the turtles!")
		return nil
	} else {
		fmt.Println("Dagger-in-Dagger inception:", meta)
	}

	host := R.client.Host()

	c = c.WithWorkdir("/work")

	sock := host.UnixSocket("/var/run/docker.sock")
	c = c.WithUnixSocket("/var/run/docker.sock", sock)

	// meta...
	file := host.File("testscript.go")
	c = c.WithFile("testscript.go", file)

	c = c.WithExec([]string{"go", "env"})
	c = c.WithExec([]string{"go", "mod", "init", "hof.io/inception"})
	c = c.WithExec([]string{"go", "mod", "tidy"})
	c = c.WithExec([]string{"go", "run", "testscript.go", fmt.Sprint(meta - 1)})

	out, err := c.Stdout(R.ctx)
	if err != nil {
		fmt.Println(out)
		return err
	}

	fmt.Println(out)
	return nil
}
