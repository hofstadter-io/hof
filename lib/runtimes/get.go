package runtimes

import (
	"fmt"
)

func RunGetFromArgs(args []string) error {
	fmt.Println("lib/runtimes.Get", args)

	return nil
}
