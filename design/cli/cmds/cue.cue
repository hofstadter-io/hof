package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// Commands drawn from Cue with a touch of hof

#DefCommand: schema.#Command & {
	TBD:   "α"
	Name:  "def"
	Usage: "def"
	Short: "print consolidated definitions"
	Long:  Short
}

#EvalCommand: schema.#Command & {
	TBD:   "α"
	Name:  "eval"
	Usage: "eval"
	Short: "print consolidated definitions"
	Long:  Short
}

#ExportCommand: schema.#Command & {
	TBD:   "α"
	Name:  "export"
	Usage: "export"
	Short: "export your data model to various formats"
	Long:  Short
}

#FormatCommand: schema.#Command & {
	TBD:   "α"
	Name:  "fmt"
	Usage: "fmt"
	Short: "formats code and files"
	Long:  Short
}

#ImportCommand: schema.#Command & {
	TBD:   "α"
	Name:  "import"
	Usage: "import"
	Short: "convert other formats and systems to hofland"
	Long:  Short
}

#TrimCommand: schema.#Command & {
	TBD:   "α"
	Name:  "trim"
	Usage: "trim"
	Short: "cleanup code, configuration, and more"
	Long:  Short
}

#VetCommand: schema.#Command & {
	TBD:   "α"
	Name:  "vet"
	Usage: "vet"
	Short: "validate data"
	Long:  Short
}
