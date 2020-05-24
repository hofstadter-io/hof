package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#ReproCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "reproduce"
	Usage: "reproduce"
	Aliases: ["repro"]
	Short: "Record, share, and replay reproducible environments and processes"
	Long:  Short
}

#UiCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "ui"
	Usage: "ui"
	Short: "Run hof's local web ui"
	Long:  Short
}

#TuiCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "tui"
	Usage: "tui"
	Short: "Run hof's terminal ui"
	Long:  Short
}

#ReplCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "repl"
	Usage: "repl"
	Short: "Run hof's local REPL"
	Long:  Short
}

#JumpCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "jump"
	Usage: "jump"
	Aliases: ["j", "leap"]
	Short: "Jumps help you do things with fewer keystrokes."
	Long:  Short
}

#DocCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "doc"
	Usage: "doc"
	Aliases: ["docs"]
	Short: "Generate and view documentation."
	Long:  Short
}
