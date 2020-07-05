package hack

import (
	"fmt"

	"github.com/hofstadter-io/hof/script"
)


func Hack(args[] string) (err error) {
	fmt.Println("Hack:", args)
	err = script.Hack(args)
	return err
}
