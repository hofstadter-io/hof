package examples

import (
	"hof.io/docs/example/gen"
)

MyGen: gen.Generator & {
	@gen(server)

	Outdir: "./output"
	// ModuleName: "hof.io/docs/example"

	// Needed because we are using the generator from within it's directory
	// Normally, users will not see or set this field
	GoModule:   "hof.io/docs/example"
	ModuleName: ""

	// We write the details in a separate file 
	"Server":    Server
	"Datamodel": Datamodel
}
