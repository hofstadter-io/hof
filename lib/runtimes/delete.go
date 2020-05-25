package runtimes

import (
	"fmt"
)

func RunDeleteFromArgs(args []string) error {
	fmt.Println("lib/runtimes.Delete", args)

	return nil
}
