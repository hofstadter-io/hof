package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func Info(args []string, rflags flags.RootPflagpole) error {
	dst := os.Getenv("DAGGER_SESSION_TOKEN")
	dcmd := "dagger run --progress tty hof env info"
	script := dcmd + " " + strings.Join(args, " ")

	envs := os.Environ()

	// incept if we are not in dagger
	if dst == "" {
		cmd := exec.Command("bash", "-c", script)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Env = envs
		err := cmd.Run()
		if err != nil {
			return err
		}

		return nil
	}

	for _, e := range envs {
		fmt.Println(e)
	}
	// R, err := prepRuntime(nil, rflags)
	// if err != nil {
	// 	return err
	// }

	return nil
}
