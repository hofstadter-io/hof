package labels

import (
	"fmt"
)

func RunInfoLabelFromArgs(args []string) error {
	fmt.Println("lib/labels.InfoLabel", args)

	return nil
}

func RunInfoLabelsetFromArgs(args []string) error {
	fmt.Println("lib/labels.InfoLabelset", args)

	return nil
}
