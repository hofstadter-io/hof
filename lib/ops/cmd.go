package ops

import (
	"fmt"
)

func RunCmdFromArgs(args []string) error {
	fmt.Println("lib/ops.Cmd", args)

	return nil
}
