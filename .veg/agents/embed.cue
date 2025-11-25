@extern(embed)

package agents

import (
  "strings"
)

_instructions: _ @embed(glob=instructions/*/*.md,type=text)
_instructions: [string]: string
instructions: [string]: [string]: string

instructions: {
  for path, content in _instructions {
    let parts = strings.Split(path, "/")
    let kind = parts[1]
    let name = strings.Split(parts[2], ".")[0]
    (kind): (name): content
  }
}
