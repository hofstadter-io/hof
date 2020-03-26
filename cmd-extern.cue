package hof

import (
	"github.com/hofstadter-io/cuemod--cli-golang/schema"

	"github.com/hofstadter-io/mvs"
)

ModCommand :: schema.Command & {
  Name:  "mod"
  Usage: "mod"
  Aliases: ["m"]
  Short:    "manage project modules"
  Long:     "Hof has mvs embedded, so you can do all the same things from this subcommand"
  Commands: mvs.CLI.Commands
}

CueCommand :: schema.Command & {
  Name:  "cue"
  Usage: "cue"
  Aliases: ["c"]
  Short: "Call a cue command"
  Long:  "Hof has cuelang embedded, so you can use hof anywhere you use cue"
  Body: """
    fmt.Println("run: cue", args)
  """
}
