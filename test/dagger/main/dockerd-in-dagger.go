package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"dagger.io/dagger"
)

const dockerVer = "docker:24"

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

	// docker-cli
	base, err := R.baseContainer()
	checkErr(err)

	// docker daemon | dind
	daemon, err := R.daemonContainer()
	checkErr(err)

	err = R.clientHack(daemon)
	checkErr(err)

	return

	// attach deemon as a service to our base image
	cntr, err := R.attachService(base, daemon)
	checkErr(err)

	go func() {
		fmt.Println("starting daemon")
		_, err = daemon.ExitCode(R.ctx)
		checkErr(err)
	}()

	fmt.Println("sleeping")
	time.Sleep(10*time.Second)

	err = R.dockerInfo(cntr)
	checkErr(err)

	err = R.dockerTest(cntr)
	checkErr(err)
}

func (R *runtime) baseContainer() (*dagger.Container, error) {

	c := R.client.Container().From(dockerVer + "-cli")
	c = c.Pipeline("base/image")
	c = c.WithExec([]string{"apk", "add", "curl"})

	return c, nil
}

func (R *runtime) daemonContainer() (*dagger.Container, error) {

	c := R.client.Container().From(dockerVer + "-dind")
	c = c.Pipeline("daemon/image")

	c = c.WithMountedCache("/tmp", R.client.CacheVolume("shared-tmp"))
	c = c.WithMountedCache("/var/lib/docker", R.client.CacheVolume("docker-cache"))
	c = c.WithExposedPort(2375)

	c = c.WithExec(
		[]string{
			"dockerd",
			"--log-level=debug",
			"--host=tcp://0.0.0.0:2375",
			"--tls=false",
		},
		dagger.ContainerWithExecOpts{ InsecureRootCapabilities: true },
	)
	return c, nil
}

func (R *runtime) attachService(c, s *dagger.Container) (*dagger.Container, error) {
	t := c.Pipeline("docker/service")
	t = t.WithEnvVariable("DOCKER_HOST", "tcp://global-dockerd:2375")
	t = t.WithServiceBinding("global-dockerd", s)
	t = t.WithMountedCache("/tmp", R.client.CacheVolume("shared-tmp"))

	return t, nil
}

func (R *runtime) dockerInfo(c *dagger.Container) error {
	t := c.Pipeline("docker/info")

	t = t.WithExec([]string{"docker", "info"})

	out, err := t.Stdout(R.ctx)
	if err != nil {
		return err
	}

	fmt.Println(out)
	return nil
}

func (R *runtime) clientHack(daemon *dagger.Container) error {

	cntr, err := R.baseContainer()
	if err != nil {
		return err
	}

	cntr, err = R.attachService(cntr, daemon)
	if err != nil {
		return err
	}

	path := "/info"
	cntr = cntr.WithExec([]string{"curl", "global-dockerd:2375" + path})

	cntr.Sync(R.ctx)

	return nil
}

func (R *runtime) dockerTest(c *dagger.Container) error {
	t := c.Pipeline("docker/test")

	t = t.WithEnvVariable("CACHE", time.Now().String())

	t = t.WithExec([]string{"docker", "pull", "hello-world"})
	t = t.WithExec([]string{"docker", "run", "hello-world"})

	t = t.WithExec([]string{"docker", "pull", "nginxdemos/hello"})
	t = t.WithExec([]string{"docker", "images"})
	t = t.WithExec([]string{"docker", "run", "-p", "4000:80", "-d", "nginxdemos/hello"})
	t = t.WithExec([]string{"curl", "global-dockerd:4000"})

	out, err := t.Stdout(R.ctx)
	if err != nil {
		return err
	}

	fmt.Println(out)

	return err
}

func (R *runtime) daemonTest(c *dagger.Container) error {
	t := c.Pipeline("daemon/test")

	t = t.WithExec([]string{"ls", "-l", "/sys/fs"})

	out, err := t.Stdout(R.ctx)
	if err != nil {
		return err
	}

	fmt.Println(out)

	return err
}
