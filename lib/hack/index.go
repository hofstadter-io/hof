package hack

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/runtime"
)


func Hack(args[] string) (err error) {
	fmt.Println("Hack:", args)

	err = runtime.GetRuntime().Print()
	if err != nil {
		// return err
	}

	err = Cueload(args)
	if err != nil {
		return err
	}

	return err
}
