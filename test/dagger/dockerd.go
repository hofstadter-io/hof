package dagger

import (
	"dagger.io/dagger"
)

const dockerVer = "docker:24"

func (R *Runtime) DockerClientContainer() (*dagger.Container, error) {

	// official docker-cli image
	c := R.Client.Container().From(dockerVer + "-cli")
	c = c.Pipeline("docker/client")

	// everyone gets a shared /tmp
	c = c.WithMountedCache("/tmp", R.Client.CacheVolume("shared-tmp"))

	return c, nil
}

// Defines an image or container that runs the dockerd daemon inside a container.
// We attach this to another container as a service, so that things running there
// can interact with the docker client, pull images, or run containers.
// We need this because hof will exec out to docker cli and then http request at the containers.
// When just mounting the docker socket, these containers were unreachable.
func (R *Runtime) DockerDaemonContainer() (*dagger.Container, error) {
	// official docker-in-docker (dind) image
	c := R.Client.Container().From(dockerVer + "-dind")
	c = c.Pipeline("docker/daemon")

	// everyone gets a shared /tmp
	c = c.WithMountedCache("/tmp", R.Client.CacheVolume("shared-tmp"))

	// this prevents overlay-in-overlay issues
	c = c.WithMountedCache("/var/lib/docker", R.Client.CacheVolume("docker-cache"))

	// docker's default port
	c = c.WithExposedPort(2375)

	// last command is the one that will be executed, dockerd in our case
	c = c.WithExec(
		[]string{
			"dockerd",
			"--log-level=warn",
			"--host=tcp://0.0.0.0:2375",
			"--tls=false",
		},
		dagger.ContainerWithExecOpts{InsecureRootCapabilities: true},
	)
	return c, nil
}

// Attaches the daemon as a dagger service to the container with the right settings.
// A single daemon can be shared among many containers this way, running separately from each of them.
func (R *Runtime) AttachDaemonAsService(c, daemon *dagger.Container) (*dagger.Container, error) {
	c = c.Pipeline("docker/attach")
	c = c.WithEnvVariable("DOCKER_HOST", "tcp://global-dockerd:2375")
	c = c.WithServiceBinding("global-dockerd", daemon)

	return c, nil
}

func (R *Runtime) AddDockerCLI(c *dagger.Container) (*dagger.Container) {
	dockerCLI := R.Client.Container().From(dockerVer).
		File("/usr/local/bin/docker")

	c = c.WithFile("/usr/local/bin/docker", dockerCLI)

	return c
}

func (R *Runtime) AddNerdctl(c *dagger.Container) (*dagger.Container) {

	ver := "1.4.0"
	//url := "https://github.com/containerd/nerdctl/releases/download/v%s/nerdctl-full-%s-linux-amd64.tar.gz"
	//tar := R.Client.HTTP(fmt.Sprintf(url, ver, ver))

	url := "https://github.com/containerd/nerdctl"
	code := R.Client.Git(url, dagger.GitOpts{ KeepGitDir: true }).
		Tag("v" + ver).
		Tree()

	c = c.Build(code)

	return c
}
