package workspace

import (
	"fmt"
)

func RunCheckoutFromArgs(args []string) error {
	fmt.Println("lib/workspace.Checkout", args)

	return nil
}
