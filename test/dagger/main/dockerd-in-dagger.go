package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"dagger.io/dagger"
)

type runtime struct {
	ctx    context.Context
	client *dagger.Client
}

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
	checkErr(err)
	defer client.Close()

	R := &runtime{
		ctx:    ctx,
		client: client,
	}

	b, err := R.baseContainer()
	checkErr(err)

	d, err := R.dockerDaemon()
	checkErr(err)

	go func() {
		fmt.Println("starting daemon")
		_, err := d.Sync(R.ctx)
		checkErr(err)
	}()

	fmt.Println("sleeping")
	time.Sleep(10*time.Second)

	c, err := R.dockerService(b, d)
	checkErr(err)

	err = R.dockerInfo(c)
	checkErr(err)

	err = R.dockerTest(c)
	checkErr(err)
}

func (R *runtime) baseContainer() (*dagger.Container, error) {

	c := R.client.Container().From("docker:24-cli")

	return c, nil
}

func (R *runtime) dockerDaemon() (*dagger.Container, error) {

	c := R.client.Container().From("docker:24-dind")
	c = c.Pipeline("docker-daemon")

	c = c.WithMountedCache("/tmp", R.client.CacheVolume("shared-tmp"))
	c = c.WithExposedPort(2375)

	c = c.WithExec(
		[]string{"dockerd", "--log-level=error", "--host=tcp://0.0.0.0:2375", "--tls=false"},
		dagger.ContainerWithExecOpts{ InsecureRootCapabilities: true },
	)
	return c, nil
}

func (R *runtime) dockerService(c, d *dagger.Container) (*dagger.Container, error) {
	t := c.Pipeline("docker-service")
	t = t.WithEnvVariable("DOCKER_HOST", "tcp://global-dockerd:2375")
	t = t.WithServiceBinding("global-dockerd", d)
	t = t.WithMountedCache("/tmp", R.client.CacheVolume("shared-tmp"))

	return t, nil
}

func (R *runtime) dockerInfo(c *dagger.Container) error {
	t := c.Pipeline("info")

	t = t.WithExec([]string{"docker", "info"})

	out, err := t.Stdout(R.ctx)
	if err != nil {
		return err
	}

	fmt.Println(out)
	return nil
}

func (R *runtime) dockerTest(c *dagger.Container) error {
	t := c.Pipeline("test")

	t = t.WithExec([]string{"docker", "pull", "nginxdemos/hello"})
	t = t.WithExec([]string{"docker", "images"})
	// t = t.WithExec([]string{"ls", "-l", "/sys/fs/cgroup"})
	t = t.WithExec([]string{"docker", "run", "-p", "4000:80", "-d", "nginxdemos/hello"})
	t = t.WithExec([]string{"curl", "localhost:4000"})

	_, err := t.Sync(R.ctx)

	return err
}
