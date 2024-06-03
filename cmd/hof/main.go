package main

import (
	"os"

	"github.com/hofstadter-io/hof/cmd/hof/cmd"
)

func main() {
	os.Setenv("CUE_EXPERIMENT", "modules=0")

	cmd.RunExit()
}
