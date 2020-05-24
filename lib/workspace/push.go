package workspace

import (
	"fmt"
)

func RunPushFromArgs(args []string) error {
	fmt.Println("lib/workspace.Push", args)

	return nil
}
