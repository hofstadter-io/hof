package labels

import (
	"fmt"
)

func RunSetLabelFromArgs(args []string) error {
	fmt.Println("lib/labels.SetLabel", args)

	return nil
}

func RunSetLabelsetFromArgs(args []string) error {
	fmt.Println("lib/labels.SetLabelset", args)

	return nil
}
