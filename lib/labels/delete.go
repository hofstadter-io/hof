package labels

import (
	"fmt"
)

func RunDeleteLabelFromArgs(args []string) error {
	fmt.Println("lib/labels.DeleteLabel", args)

	return nil
}

func RunDeleteLabelsetFromArgs(args []string) error {
	fmt.Println("lib/labels.DeleteLabelset", args)

	return nil
}
