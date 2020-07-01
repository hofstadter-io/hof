package runtime

import (
	"fmt"
)

func (ts *Script) CmdBilly(neg int, args []string) {
	fmt.Println("billy!", neg, args)

}
