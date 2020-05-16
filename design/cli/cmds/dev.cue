package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#DevCommand: schema.#Command & {
  Name:    "dev"
  Usage:   "dev"
  Short:   "run hof's local dev setup"
  Long:    Short
}

#UiCommand: schema.#Command & {
  Name:    "ui"
  Usage:   "ui"
  Short:   "run hof's local web ui"
  Long:    Short
}

#ReplCommand: schema.#Command & {
  Name:    "repl"
  Usage:   "repl"
  Short:   "run hof's REPL system"
  Long:    Short
}
