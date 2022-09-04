package yagu

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func Bash(script, workdir string) (string, error) {
	cmd := exec.Command("bash", "-p", "-c", script)
	if workdir != "" {
		cmd.Dir = workdir
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func BashTmpScriptWithArgs(script string, args []string) (string, error) {

	// make sure we have our temp dir
	err := Mkdir(".hof")
	if err != nil {
		return "", err
	}

	// create temp file
	f, err := ioutil.TempFile(".hof", "tmp-jumps-*")
	if err != nil {
		return "", err
	}
	// and cleanup for sure
	defer os.Remove(f.Name())

	// write script
	_, err = fmt.Fprintln(f, script)
	if err != nil {
		return "", err
	}

	// close
	err = f.Close()
	if err != nil {
		return "", err
	}

	// make exec'd
	err = os.Chmod(f.Name(), 0755)
	if err != nil {
		return "", err
	}

	// prepare args to bash invocation as a "/path/to/scripts.sh"
	doargs := append([]string{f.Name()}, args...)
	dorun := strings.Join(doargs, " ")

	// run script and return output / status
	cmd := exec.Command("bash", "-p", "-c", dorun)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func Exec(args []string) (string, error) {
	command, rest := args[0], args[1:]
	cmd := exec.Command(command, rest...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
