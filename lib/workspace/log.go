package workspace

import (
	"fmt"
)

func RunLogFromArgs(args []string) error {
	fmt.Println("lib/workspace.Log", args)

	return nil
}
