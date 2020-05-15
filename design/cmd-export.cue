package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#ExportCommand: schema.#Command & {
  Name:    "export"
  Usage:   "export"
  Short:   "export your data model to various formats"
  Long:    Short
},

