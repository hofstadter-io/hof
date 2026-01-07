package env

import (
	"github.com/hofstadter-io/hof/schemas"
)


// supported formats for sbom generators

// common fields for sbom generators, used internally
sbomCommon: {
	name?:    string
	path:    string
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
	format:  or(["cue", "json", "yaml", "toml"])

  sbomCommon 
}

// generate the Dagger representation
// is a: *dagger.File with JSON content
// data: any *dagger.Object (env.Ref, i.e. #Things)
// hmmm, can we reverse this one?
// update, can't seem to get anything reasonable out of dagger for sbom, misleading function names / what they return, it's all internal ids to ephemeral object, not actual sbom material
// #DaggerSBOM: Ref & {
// 	schemas.Hof
// 	#hof: env: {
// 		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
// 		kind: "daggerSBOM"
// 	}
// 	$kind:   "#daggerSBOM"

//   sbomCommon 
// }

// TODO, sigstore/cosign stuff
