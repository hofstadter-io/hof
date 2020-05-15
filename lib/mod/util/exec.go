package util

import (
	"os/exec"
)

func Bash(script string) (string, error) {
	cmd := exec.Command("bash", "-c", script)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func Exec(args []string) (string, error) {
	command, rest := args[0], args[1:]
	cmd := exec.Command(command, rest...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
