package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// Kubernetes inspired commands (maybe some hyper-cloud too)

#LabelCommand: schema.#Command & {
  Name:  "label"
  Usage: "label"
  Short: "manage resource labels"
	Long: Short
}

#CreateCommand: schema.#Command & {
  Name:  "create"
  Usage: "create"
  Short: "create resources"
	Long: Short
}

#ApplyCommand: schema.#Command & {
  Name:  "apply"
  Usage: "apply"
  Short: "apply resource configuration"
	Long: Short
}

#GetCommand: schema.#Command & {
  Name:  "get"
  Usage: "get"
  Short: "find and display resources"
	Long: Short
}

#DeleteCommand: schema.#Command & {
  Name:  "delete"
  Usage: "delete"
  Short: "delete resources"
	Long: Short
}

