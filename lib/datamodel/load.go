package datamodel

import (
	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/cuetils"
)


func LoadDatamodel(entrypoints []string) (cue.Value, error) {
	crt, err := cuetils.CueRuntimeFromEntrypointsAndFlags(entrypoints)
	if err != nil {
		return cue.Value{}, err
	}

	val := crt.CueValue

	return val, nil
}
