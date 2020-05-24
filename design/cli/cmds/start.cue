package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#SetupCommand: schema.#Command & {
	TBD:   "α"
	Name:  "setup"
	Usage: "setup"
	Short: "setup the hof tool"
	Long:  Short
}

#CloneCommand: schema.#Command & {
	TBD:   "α"
	Name:  "clone"
	Usage: "clone"
	Short: "clone a workspace or repository into a new directory"
	Long:  Short
}

#InitCommand: schema.#Command & {
	TBD:   "α"
	Name:  "init"
	Usage: "init"
	Short: "create an empty workspace or initialize an existing directory to one"
	Long:  Short
}
