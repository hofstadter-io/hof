package advanced

import (
	"github.com/hofstadter-io/hof/catalogs/env/bases"
	"github.com/hofstadter-io/hof/schemas/env"
)

without: env.#Container & {
  @env(without)
  from: bases.rocky8.minimal
  steps: [
    env.WithoutFile & { path: "/usr/bin/curl" }
  ]
}