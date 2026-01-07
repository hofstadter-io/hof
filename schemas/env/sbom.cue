package env

import (
	"github.com/hofstadter-io/hof/schemas"
)


// supported formats for sbom generators
sbomFormats: ["cue", "json", "yaml", "toml"]

// common fields for sbom generators, used internally
sbomCommon: {
	name?:    string
	path:    string
	format:  or(sbomFormats)
	data:    _
}

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

  sbomCommon 
}

// generate the Dagger representation
// is a: *dagger.File with JSON content
// data: any *dagger.Object (env.Ref, i.e. #Things)
// hmmm, can we reverse this one?
#DaggerSBOM: Ref & {
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "daggerSBOM"
	}
	$kind:   "#daggerSBOM"

  sbomCommon 
}

// TODO, sigstore/cosign stuff
