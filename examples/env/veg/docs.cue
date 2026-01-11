@experiment(aliasv2)
package veg

import (

	"github.com/hofstadter-io/hof/schemas/env"
	// "strings"
)

let root = self

docs: {

	ctr: env.#Container & {
		@env(docs-ctr)
		from: root.ctr.dev
		steps: [
			env.Dir & {source: root.src.docs},
		]
	}

}
