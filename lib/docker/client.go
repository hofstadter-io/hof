package docker

import (
	"context"
	"fmt"
	"time"

	credClient "github.com/docker/docker-credential-helpers/client"
	"github.com/docker/docker/client"
)

var dockerClient *client.Client

func InitDockerClient() (err error) {
	dockerClient, err = client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
		client.WithHostFromEnv(),
	)
	if err != nil {
		return fmt.Errorf("error: hof fmt requires docker\n%w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := dockerClient.Ping(ctx); err != nil {
		return fmt.Errorf("docker ping: %w", err)
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
