package env

import (
	"github.com/hofstadter-io/hof/schemas"
)


// generate the CUE representation
// is a: *dagger.File with format:[cue,json,yaml,toml] content
// data: any CUE value
// hmmm, can we reverse this one?
#CuefigSBOM: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "cuefigSBOM"
	}
	$kind:   "#cuefigSBOM"
	format:  or(["cue", "json", "yaml", "toml"])

	name?:    string
	path:    string
	data:    _
}

// TODO, sigstore/cosign stuff
