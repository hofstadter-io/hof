package labels

import (
	"fmt"
)

func RunGetLabelFromArgs(args []string) error {
	fmt.Println("lib/labels.GetLabel", args)

	return nil
}

func RunGetLabelsetFromArgs(args []string) error {
	fmt.Println("lib/labels.GetLabelset", args)

	return nil
}
