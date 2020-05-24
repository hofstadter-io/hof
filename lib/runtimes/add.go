package runtimes

import (
	"fmt"
)

func RunAddFromArgs(args []string) error {
	fmt.Println("lib/runtimes.Add", args)

	return nil
}
