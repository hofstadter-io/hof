package cuetils

import (
	"cuelang.org/go/cue/load"
)

// CueRuntimeFromArgs builds up a CueRuntime
//  by processing the args passed in
func CueRuntimeFromArgs(args []string) (crt *CueRuntime, err error) {
	crt = &CueRuntime{
		Entrypoints: args,
		CueConfig: &load.Config{
			ModuleRoot: "",
			Module: "",
			Package: "",
			Dir: "",
			BuildTags: []string{},
			Tests: false,
			Tools: false,
			DataFiles: false,
			Overlay: map[string]load.Source{},
		},
	}

	err = crt.Load()

	return crt, err
}

// CueRuntimeFromArgsAndFlags builds up a CueRuntime
//  by processing the args passed in AND the current flag values
func CueRuntimeFromArgsAndFlags(args []string) (crt *CueRuntime, err error) {
	crt = &CueRuntime{
		Entrypoints: args,
	}

	// XXX TODO XXX
	// Buildup out arg to load.Instances second arg
	// Add this configuration to our runtime struct

	err = crt.Load()

	return crt, err
}

