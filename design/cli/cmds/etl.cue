package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#EtlCommand: schema.#Command & {
  Name:  "etl"
  Usage: "etl"
  Short: "extract, transform, and load data"
	Long: Short
}

