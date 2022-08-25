package create

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func Create(args []string, rootflags flags.RootPflagpole, cmdflags flags.CreateFlagpole) error {

	fmt.Println("Create:", args, cmdflags)

	return nil
}
