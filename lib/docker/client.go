package docker

import (
	"fmt"

	"github.com/docker/docker/client"
	credClient "github.com/docker/docker-credential-helpers/client"
)

var dockerClient *client.Client

func InitDockerClient() (err error) {
	dockerClient, err = client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return fmt.Errorf("error: hof fmt requires docker\n%w", err)
	}
	return nil
}

func Hack() (err error) {
	p := credClient.NewShellProgramFunc("docker-credential-gcloud")
	creds, err := credClient.Get(p, "https://us.gcr.io")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Got credentials for user `%s` in `%s`\n", creds.Username, creds.ServerURL)

	return nil
}
