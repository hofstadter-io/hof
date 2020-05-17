package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#UiCommand: schema.#Command & {
  Name:    "ui"
  Usage:   "ui"
  Short:   "Run hof's local web ui"
  Long:    Short
}

#TuiCommand: schema.#Command & {
  Name:    "tui"
  Usage:   "tui"
  Short:   "Run hof's terminal ui"
  Long:    Short
}

#ReplCommand: schema.#Command & {
  Name:    "repl"
  Usage:   "repl"
  Short:   "Run hof's local REPL"
  Long:    Short
}
