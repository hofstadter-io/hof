package container

import (
	"fmt"
	"os"
	"os/exec"
)

var client Client

const (
	envRuntime = "HOF_CONTAINER_RUNTIME"
)

var rt Runtime

type Client struct {
	runtimePath string
}

func InitClient() error {
	var (
		rb       RuntimeBinary
		binaries = []RuntimeBinary{
			RuntimeBinary(os.Getenv(envRuntime)),
			RuntimeBinaryDocker,
			RuntimeBinaryPodman,
			RuntimeBinaryNerdctl,
		}
	)

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
