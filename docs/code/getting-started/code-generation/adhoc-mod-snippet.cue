package example

import (
	"github.com/hofstadter-io/hof/schema/gen"
}

foo: gen.#Generator: {
	@gen(foo)

	// input data
	In: _

	Out: [
		// list of files to generate
	]

	// other fields filled by hof
}
