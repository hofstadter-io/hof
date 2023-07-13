package container

type docker struct {
	runtime
}

func newDocker() docker {
	return docker{
		runtime: newRuntime(RuntimeBinaryDocker),
	}
}
