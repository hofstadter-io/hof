package container

type nerdctl struct {
	runtime
}

func newNerdctl() nerdctl {
	return nerdctl{
		runtime: newRuntime(RuntimeBinaryNerdctl),
	}
}
