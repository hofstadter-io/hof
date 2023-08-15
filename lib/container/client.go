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

type Client struct {
	runtimePath string
}

func InitClient() error {
	urt := os.Getenv(envRuntime)

	// short-circuit if none is explicitly set
	if urt == "none" {
		rt = newNone()
		return nil
	}

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
	case "none":
		rt = newNone()
	default:
		fmt.Println("failed to find any container runtimes %s in PATH", binaries)
		fmt.Println("set HOF_CONTAINER_RUNTIME=none to disable this message")
		rt = newNone()
	}

	return nil
}
