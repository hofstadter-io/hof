package script

import (
	"fmt"

	"github.com/hofstadter-io/hof/script/runtime"
)

func Hack(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please supply at least a file/path.hls")
	}

	err := runtime.RunScript(args)
	if err != nil {
		return fmt.Errorf("in script.Hack")
	}

	return nil
}


