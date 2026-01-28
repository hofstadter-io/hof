package agentic

import (
	"github.com/hofstadter-io/hof/schemas"
)

Tool: {
  schema.Hof

  #hof: agentic: { root: true, kind: "tool" }

  name: string

}