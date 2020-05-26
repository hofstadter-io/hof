package cuetils

import (
	"cuelang.org/go/cue/load"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

// CueRuntimeFromArgs builds up a CueRuntime
//  by processing the args passed in
func CueRuntimeFromEntrypoints(entrypoints []string) (crt *CueRuntime, err error) {
	crt = &CueRuntime{
		Entrypoints: entrypoints,
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
func CueRuntimeFromEntrypointsAndFlags(entrypoints []string) (crt *CueRuntime, err error) {
	cfg := &load.Config{
		ModuleRoot: "",
		Module: "",
		Package: "",
		Dir: "",
		BuildTags: []string{},
		Tests: false,
		Tools: false,
		DataFiles: false,
		Overlay: map[string]load.Source{},
	}

	// package?
	if flags.RootPackagePflag != "" {
		cfg.Package = flags.RootPackagePflag
	}

	crt = &CueRuntime{
		Entrypoints: entrypoints,
		CueConfig: cfg,
	}

	err = crt.Load()

	return crt, err
}

