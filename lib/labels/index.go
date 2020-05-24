package labels

import (
	"fmt"
)

// Should be able to query select by labels and then apply one more more lables
// So basically "Get" with label changes, will need to expand to some CRUD ops
// almost their own thing, but want some real flexibility here
func LabelOperationFromArgs(args []string) error {
	fmt.Println("lib/resources.Label", args)

	return nil
}
