package runtimes

import (
	"fmt"
)

func RunUninstallFromArgs(args []string) error {
	fmt.Println("lib/runtimes.Uninstall", args)

	return nil
}
