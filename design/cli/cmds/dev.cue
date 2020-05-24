package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#ReproCommand: schema.#Command & {
	TBD:   "Ø"
	Name:  "reproduce"
	Usage: "reproduce"
	Aliases: ["repro"]
	Short: "Record, share, and replay reproducible environments and processes"
	Long:  Short
}

#UiCommand: schema.#Command & {
	TBD:   "Ø"
	Name:  "ui"
	Usage: "ui"
	Short: "Run hof's local web ui"
	Long:  Short
}

#TuiCommand: schema.#Command & {
	TBD:   "Ø"
	Name:  "tui"
	Usage: "tui"
	Short: "Run hof's terminal ui"
	Long:  Short
}

#ReplCommand: schema.#Command & {
	TBD:   "Ø"
	Name:  "repl"
	Usage: "repl"
	Short: "Run hof's local REPL"
	Long:  Short
}

#JumpCommand: schema.#Command & {
	TBD:  "α"
	Name:  "jump"
	Usage: "jump"
	Aliases: ["j", "leap"]
	Short: "Jumps help you do things with fewer keystrokes."
	Long:  Short
}

#DocCommand: schema.#Command & {
	TBD:   "Ø"
	Name:  "doc"
	Usage: "doc"
	Aliases: ["docs"]
	Short: "Generate and view documentation."
	Long:  Short
}
