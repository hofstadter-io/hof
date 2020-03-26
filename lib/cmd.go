package lib

import (
	"fmt"
)

func Cmd(flags, args []string, mode string) (string, error) {
	fmt.Println("Cmd", flags, args)

	// ... Cue SDK, simulate cue eval / export
	// Pick out and export anything starting with Gen
	//   if we can determine if struct or list and loop appropriately

	// see if we can parse and introspect *_tool.cue files

	return "not implemented", nil
}
