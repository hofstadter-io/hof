package cmd

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func Run(args []string, rflags flags.RootPflagpole, cflags flags.AgentFlagpole) error {

	fmt.Println("Agent!")

	return nil
}
