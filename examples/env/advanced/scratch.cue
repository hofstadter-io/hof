package advanced

import (
	"github.com/hofstadter-io/hof/catalogs/env/bases"
	// "github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

scratch: {
  prep: env.#Container & {
    @env(scratch-prep)
    from: bases.alpine3.minimal
  }
  image: env.#Container & {
    @env(scratch-image)
    from: "scratch"
  }
}