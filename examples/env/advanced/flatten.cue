package advanced

import (
	"github.com/hofstadter-io/hof/catalogs/env/bases"
	"github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

// flatten demonstration
flatten: {

  // baseline
  base: env.#Container & {
    @env(flatten/base)
    from: bases.rocky8.minimal
  }

  // w/o flatten optimization
  orig: env.#Container & {
    @env(flatten/orig)
    from: base
    steps: [
      utils.dnf.upgrade,
    ]
  }

  // with flatten optimization
  fixd: env.#Flatten & { #orig: orig, @env(flatten/fixd) }

}