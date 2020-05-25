package labels

import (
	"fmt"
)

func RunCreateLabelFromArgs(args []string) error {
	fmt.Println("lib/labels.CreateLabel", args)

	return nil
}

func RunCreateLabelsetFromArgs(args []string) error {
	fmt.Println("lib/labels.CreateLabelset", args)

	return nil
}
