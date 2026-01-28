package agentic

import (
	"github.com/hofstadter-io/hof/schemas"
)

Model: {
  schema.Hof

  #hof: agentic: { root: true, kind: "model" }

  name: string
  id:   string

  features?: _
  prices?: _
}