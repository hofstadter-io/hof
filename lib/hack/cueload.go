package hack

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/runtime"
)

func Cueload(args []string) (err error) {
	if len(args) == 0 {
		args = []string{"./"}
	}
	// fmt.Println("CueLoad:", args)

	crt, err := runtime.CueRuntimeFromArgsAndFlags(args)
	if err != nil {
		fmt.Println("Error")
		return err
	}

	// fmt.Println("errs:", crt.CueErrors)

	crt.PrintValue()

	return err
}
