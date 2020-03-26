package hof

import (
	"github.com/hofstadter-io/cuemod--cli-golang/schema"
)

GenCommand :: schema.Command & {
  Name:  "gen"
  Usage: "gen [files...]"
  Aliases: ["g"]
  Short: "generate code, data, and config"
  Long: """
    generate all the things, from code to data to config...
  """

  Imports: [
    {Path: "fmt"},
    {Path: "os"},
    {Path: CLI.Package + "/lib"},
  ]

  Body: """
    msg, err := lib.Gen(args, []string{}, "")
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
    fmt.Println(msg)
  """
}
