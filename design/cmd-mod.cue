package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"

	"github.com/hofstadter-io/mvs"
)

#ModCommand: schema.#Command & {
  Name:  "mod"
  Usage: "mod"
  Aliases: ["m"]
  Short:    "manage project modules"
  Long:     "Hof has mvs embedded, so you can do all the same things from this subcommand"
  Commands: mvs.#CLI.Commands
}

