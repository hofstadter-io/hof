package workspace

import (
	"fmt"
)

func RunCloneFromArgs(module, name string) error {
	fmt.Println("lib/workspace.Clone", module, name)

	return nil
}
