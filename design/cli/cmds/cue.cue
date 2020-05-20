package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// Commands drawn from Cue with a touch of hof

#DefCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "def"
	Usage: "def"
	Short: "print consolidated definitions"
	Long:  Short
}

#EvalCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "eval"
	Usage: "eval"
	Short: "print consolidated definitions"
	Long:  Short
}

#ExportCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "export"
	Usage: "export"
	Short: "export your data model to various formats"
	Long:  Short
}

#FormatCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "fmt"
	Usage: "fmt"
	Short: "formats code and files"
	Long:  Short
}

#ImportCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "import"
	Usage: "import"
	Short: "convert other formats and systems to hofland"
	Long:  Short
}

#TrimCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "trim"
	Usage: "trim"
	Short: "cleanup code, configuration, and more"
	Long:  Short
}

#VetCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "vet"
	Usage: "vet"
	Short: "validate data"
	Long:  Short
}
