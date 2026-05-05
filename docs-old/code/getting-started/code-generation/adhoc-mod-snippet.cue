package example

import (
	"github.com/hofstadter-io/hof/schema/gen"
}

foo: gen.Generator & {
	@gen(foo)

	// input data
	In: _

	// normally when writing generators as code
	// you add the CUE to turn In -> Out
	//   - provide project specific config and flags
	//   - dynamically decide what files to generate
	//   - craft schemas and DSLs to create anything

	// list of files to generate
	Out: [...gen.File]

	// other fields filled by hof when you turn adhoc -> reusable
}
