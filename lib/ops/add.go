package ops

import (
	"fmt"
)

func RunAddFromArgs(args []string) error {
	fmt.Println("lib/ops.Add", args)

	return nil
}
