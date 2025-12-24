package container

type docker struct {
	runtime
}

func newDocker() docker {
	return docker{
		runtime: newRuntime(RuntimeBinaryDocker),
	}
}

type nerdctl struct {
	runtime
}

func newNerdctl() nerdctl {
	return nerdctl{
		runtime: newRuntime(RuntimeBinaryNerdctl),
	}
}

type podman struct {
	runtime
}

func newPodman() podman {
	return podman{
		runtime: newRuntime(RuntimeBinaryPodman),
	}
}

type none struct {
	runtime
}

func newNone() none {
	return none{
		runtime: newRuntime(RuntimeBinaryNone),
	}
}
