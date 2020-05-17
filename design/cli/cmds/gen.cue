package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#GenCommand: schema.#Command & {
  Name:  "gen"
  Usage: "gen [files...]"
  Aliases: ["g"]
  Short: "generate code, data, and config"
  Long: """
    generate all the things, from code to data to config...
  """

  Flags: [...schema.#Flag] & [
    {
      Name:    "generator"
      Type:    "[]string"
      Default: "nil"
      Help:    "Generators to run"
      Long:    "generator"
      Short:   "g"
    },
  ]

}
