package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// Kubernetes inspired commands (maybe some hyper-cloud too)

#LabelCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "label"
	Usage: "label"
	Short: "manage resource labels"
	Long:  Short
}

#CreateCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "create"
	Usage: "create"
	Short: "create resources"
	Long:  Short
}

#ApplyCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "apply"
	Usage: "apply"
	Short: "apply resource configuration"
	Long:  Short
}

#GetCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "get"
	Usage: "get"
	Short: "find and display resources"
	Long:  Short
}

#DeleteCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "delete"
	Usage: "delete"
	Short: "delete resources"
	Long:  Short
}
