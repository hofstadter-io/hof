package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#InitCommand: schema.#Command & {
  Name:    "init"
  Usage:   "init"
  Short:   "init the current directory for dma usage."
  Long:    Short
}

