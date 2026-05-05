//go:build !unix
// +build !unix

package incept

import "os/exec"

func ensureChildProcessesAreKilled(cmd *exec.Cmd) {
}
