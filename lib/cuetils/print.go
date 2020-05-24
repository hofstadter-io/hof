package cuetils

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"
)


func (CRT *CueRuntime) PrintValue() error {
	// Get top level struct from cuelang
	S, err := CRT.CueValue.Struct()
	if err != nil {
		return err
	}

	iter := S.Fields()
	for iter.Next() {

		label := iter.Label()
		value := iter.Value()
		fmt.Println("  -", label, value)
		for attrKey, attrVal := range value.Attributes() {
			fmt.Println("  --", attrKey)
			for i := 0; i < 5; i++ {
				str, err := attrVal.String(i)
				if err != nil {
					break
				}
				fmt.Println("  ---", str)
			}
		}
	}

	return nil
}
