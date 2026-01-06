package veg

import (
  "github.com/hofstadter-io/hof/schemas/env"
)

hack: {


  diskUsage: env.#Container & {
    @env(hack-diskUsage)
    from: "ghcr.io/hofstadter-io/veg-hof:v0.7.0-alpha.1"
    steps: [
			env.Mount & {path: "/work", source: src.code},
    ]
  }
}