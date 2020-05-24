package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// cue run + hof extra
#CmdCommand: schema.#Command & {
	TBD:   "α"
	Name:  "cmd"
	Usage: "cmd [flags] [cmd] [args]"
	Short: "run commands from the scripting layer and your _tool.cue files"
	Long:  Short
}
