package workspace

import (
	"fmt"
)

func RunIncludeFromArgs(args []string) error {
	fmt.Println("lib/workspace.Include", args)

	return nil
}
