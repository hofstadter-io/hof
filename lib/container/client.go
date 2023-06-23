package container

import (
	"fmt"
<<<<<<< HEAD
	"os/exec"

	credClient "github.com/docker/docker-credential-helpers/client"
)

var (
	client   Client
	binaries = []RuntimeBinary{
		RuntimeBinaryNerdctl,
		RuntimeBinaryPodman,
		RuntimeBinaryDocker,
	}
)

var rt Runtime

type Client struct {
	runtimePath string
}

func InitClient() error {
	var rb RuntimeBinary

	for _, b := range binaries {
		if _, err := exec.LookPath(string(b)); err == nil {
			rb = b
			break
		}
	}

	switch rb {
	case RuntimeBinaryNerdctl:
		rt = newNerdctl()
	case RuntimeBinaryPodman:
		rt = newPodman()
	case RuntimeBinaryDocker:
		rt = newDocker()
	default:
		return fmt.Errorf("failed to find any of %s in PATH", binaries)
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
