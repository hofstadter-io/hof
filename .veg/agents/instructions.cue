@extern(embed)

package agents

import (
  "strings"
)

// todo, something like this under .veg will be default
// we also need to combine home & project agentic stuff

instructionsDir: ".veg/agents/instructions/**/*.md"

_i1: _ @embed(glob=instructions/*/*.md,type=text)
_i2: _ @embed(glob=instructions/*/*/*.md,type=text)
// _i3: _ @embed(glob=instructions/*/*/*/*.md,type=text)

instructions: {
  for path, content in _i1 { (strings.TrimPrefix(path, "instructions/")): content }
  for path, content in _i2 { (strings.TrimPrefix(path, "instructions/")): content }
  // for path, content in _i3 { (strings.TrimPrefix(path, "instructions/")): content }
}
