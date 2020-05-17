package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// Commands drawn from Cue with a touch of hof

#DefCommand: schema.#Command & {
  Name:  "def"
  Usage: "def"
  Short: "print consolidated definitions"
	Long: Short
}

#EvalCommand: schema.#Command & {
  Name:  "eval"
  Usage: "eval"
  Short: "print consolidated definitions"
	Long: Short
}

#ExportCommand: schema.#Command & {
  Name:    "export"
  Usage:   "export"
  Short:   "export your data model to various formats"
  Long:    Short
}

#FormatCommand: schema.#Command & {
  Name:  "fmt"
  Usage: "fmt"
  Short: "formats code and files"
	Long: Short
}

#ImportCommand: schema.#Command & {
  Name:    "import"
  Usage:   "import"
  Short:   "convert other formats and systems to hofland"
	Long: Short
}

#TrimCommand: schema.#Command & {
  Name:  "trim"
  Usage: "trim"
  Short: "cleanup code, configuration, and more"
	Long: Short
}

#VetCommand: schema.#Command & {
  Name:  "vet"
  Usage: "vet"
  Short: "validate data"
	Long: Short
}
