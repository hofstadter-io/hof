package hack

import (
	"fmt"

	// "github.com/hofstadter-io/hof/lib/runtime"
)


func Hack(args[] string) (err error) {
	fmt.Println("Hack:", args)

	/*
	rt := runtime.GetRuntime()

	fmt.Println("config:", rt.ConfigType)
	runtime.GetRuntime().PrintConfig()

	fmt.Println("secret:", rt.SecretType)
	runtime.GetRuntime().PrintSecret()

	if len(args) == 1 {
		fmt.Println("==========")
		runtime.GetRuntime().ConfigGet(args[0])
	}

	if len(args) == 2 {
		fmt.Println("==========")
		runtime.GetRuntime().ConfigSet(args[0], args[1])
	}
	err = Cueload(args)
	if err != nil {
		return err
	}
	*/

	return err
}
