package examples

import (
	"hof.io/docs/example/gen"
)

Server: gen.#Generator & {
	@gen(server)

	Outdir:   "./output"

	// We write the design in a separate file 
	Server: ServerDesign

	// Needed because we are using the generator from within it's directory
	// Users who import your generator as a module will not need to set this
	// It's a temporary requirement until CUE supports embedded files
	PackageName: ""
}
