package container

type none struct {
	runtime
}

func newNone() none {
	return none{
		runtime: newRuntime(RuntimeBinaryNone),
	}
}
