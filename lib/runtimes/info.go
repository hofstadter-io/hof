package runtimes

import (
	"fmt"
)

func RunInfoFromArgs(args []string) error {
	fmt.Println("lib/runtimes.Info", args)

	return nil
}
