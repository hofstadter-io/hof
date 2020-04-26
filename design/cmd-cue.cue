package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#CueCommand: schema.#Command & {
  Name:  "cue"
  Usage: "cue"
  Aliases: ["c"]
  Short: "Call a cue command"
  Long:  "Hof has cuelang embedded, so you can use hof anywhere you use cue"

}
