package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#BuildCommand: schema.#Command & {
	TBD:   "Ø"
	Name:  "build"
	Usage: "build [flags] [cmd] [args]"
	Short: "Build assets for modules and generated output"
	Long:  Short
}
