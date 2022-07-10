package demo

import (
	"github.com/hofstadter-io/hof/schema/gen"
)

// This is example usage of your generator
DemoExample: #DemoGenerator & {
	@gen(demo)

	// inputs to the generator
	"users": users
	"data":  data

	// other settings
	Diff3:  false
	Outdir: "./"

	// required by examples inside the same module
	// your users do not set or see this field
	PackageName: ""
}

// This is your reusable generator module
//
#DemoGenerator: gen.#Generator & {

	//
	// user input fields
	//

	// this is the interface for this generator module
	// typically you enforce schema(s) here
	users: _
	data:  _

	//
	// Internal Fields
	//

	// This is the global input data the templates will see
	// You can reshape and transform the user inputs
	// While we put it under internal, you can expose In
	In: {
		// if you want to user your input data
		// add top-level fields from your
		// CUE entrypoints here, adjusting as needed
		// Since you made this a module for others,
		// it won't output until this field is filled

		"users": users
		"data":  data

		...
	}

	// required for hof CUE modules to work
	// your users do not set or see this field
	PackageName: string | *"hof.io/demo"

	// Templates: [{Globs: ["./templates/**/*"], TrimPrefix: "./templates/"}]
	Templates: [ {Globs: [ "min.txt", "ext.txt"]}]

	// Partials: [#Templates & {Globs: ["./partials/**/*"], TrimPrefix: "./partials/"}]
	Partials: []

	// The final list of files for hof to generate
	Out: [...gen.#File] & [
		t_0,
		t_1,

	]

	// These are the -T mappings
	t_0: {
		TemplatePath: "min.txt"
		Filepath:     "out.txt"
	}
	t_1: {
		TemplatePath: "ext.txt"
		Filepath:     "users.txt"
	}

	// so your users can build on this
	...
}
