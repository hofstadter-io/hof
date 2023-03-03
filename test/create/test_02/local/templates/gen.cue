package main

import (
	"github.com/hofstadter-io/hof/gen"
)

"{{ .name }}": gen.#Generator & {
	@gen({{ .name }})

		Out: []
}
