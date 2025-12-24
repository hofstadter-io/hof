package veg

import (
  "github.com/hofstadter-io/hof/lib/env/devex"
)

_reg: "host.docker.internal:5000"

environs: [n=string]: { name: string | *n }
environs: {
  for k, env in devex.veg {
    (k): {
      name: env.name
      description: env.description | "\(env.name) image"
      spec: {
        from: "\(_reg)/\(env.name):local"
      }
    }
  }
}