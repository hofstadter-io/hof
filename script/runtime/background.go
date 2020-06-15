package runtime

import "os/exec"

type backgroundCmd struct {
	cmd  *exec.Cmd
	wait <-chan struct{}
	neg  int // if true, cmd should fail
}

