package design

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

  Imports: [
    {Path: #Module + "/lib"},
  ]

  Body: """
    return lib.Gen([]string{}, []string{}, "")
  """
}
