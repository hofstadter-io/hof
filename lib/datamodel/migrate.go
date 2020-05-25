package datamodel

import (
	"fmt"
)

func RunMigrateFromArgs(args []string) error {
	fmt.Println("lib/datamodel.Migrate", args)

	return nil
}
