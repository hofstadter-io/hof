package examples

import (
	"hof.io/docs/example/gen"
)

// A runnable generator (@gen(<name>))
Generator: gen.Generator & {
	@gen(server)

	Outdir: "./output"

	// The full design lives in a separate file 
	"Server": Server

	// Needed because we are using the generator from within it's directory
	// Normally, users will not see or set this field
	ModuleName: ""
}
