package runtime

import (
	"fmt"
)

func (RT *Runtime) GetStdout() string {
	if RT.lastcmd == nil {
		return ""
	}
	return RT.lastcmd.Result().Stdout.(fmt.Stringer).String()
}

func (RT *Runtime) GetStderr() string {
	if RT.lastcmd == nil {
		return ""
	}
	return RT.lastcmd.Result().Stderr.(fmt.Stringer).String()
}

func (RT *Runtime) GetStatus() int {
	return RT.status
}
