package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// Kubernetes inspired commands (maybe some hyper-cloud too)

#InfoCommand: schema.#Command & {
	TBD:  "α"
	Name:  "info"
	Usage: "info"
	Short: "print information about known resources"
	Long:  Short
}

#CreateCommand: schema.#Command & {
	TBD:  "α"
	Name:  "create"
	Usage: "create"
	Short: "create resources"
	Long:  Short
}

#ApplyCommand: schema.#Command & {
	TBD:  "α"
	Name:  "apply"
	Usage: "apply"
	Short: "apply resource configuration"
	Long:  Short
}

#GetCommand: schema.#Command & {
	TBD:  "α"
	Name:  "get"
	Usage: "get"
	Short: "find and display resources"
	Long:  Short
}

#EditCommand: schema.#Command & {
	TBD:  "α"
	Name:  "edit"
	Usage: "edit"
	Short: "edit resources"
	Long:  Short
}

#DeleteCommand: schema.#Command & {
	TBD:  "α"
	Name:  "delete"
	Usage: "delete"
	Short: "delete resources"
	Long:  Short
}
