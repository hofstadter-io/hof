package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// cue get + hof add
#AddCommand: schema.#Command & {
	TBD:  "α"
	Name:  "add"
	Usage: "add"
	Short: "add dependencies and new components to the current module or workspace"
	Long:  Short
}
