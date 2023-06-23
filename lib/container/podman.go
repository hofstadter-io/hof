package container

type podman struct {
	runtime
}

func newPodman() podman {
	return podman{
		runtime: newRuntime(RuntimeBinaryPodman),
	}
}
