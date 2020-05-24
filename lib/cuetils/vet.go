package cuetils

import (
	"fmt"
)

func RunVetFromArgs(args []string) error {
	fmt.Println("vet", args)

	return nil
}
