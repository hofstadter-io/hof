package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

CueCommand :: schema.Command & {
  Name:  "cue"
  Usage: "cue"
  Aliases: ["c"]
  Short: "Call a cue command"
  Long:  "Hof has cuelang embedded, so you can use hof anywhere you use cue"
  Body: """
    fmt.Println("run: cue", args)
    fmt.Println("not implemented, will have to mirror cue cli command structure, args, and flags in the design here, or break out into another repo")
  """
  Imports: [
    {Path: "fmt"},
  ]
}
