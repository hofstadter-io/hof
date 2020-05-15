package util

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"
)

var (
	// Global cue runtime
	CueRuntime cue.Runtime
)

func PrintCueInstance(i *cue.Instance) error {
	bytes, err := format.Node(i.Value().Syntax())
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}
