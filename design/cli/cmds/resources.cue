package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// Kubernetes inspired commands (maybe some hyper-cloud too)

#InfoCommand: schema.#Command & {
	TBD:   "α"
	Name:  "info"
	Usage: "info"
	Short: "print information about known resources"
	Long:  Short
	Flags: [{
		Name:    "builtin"
		Type:    "bool"
		Default: "false"
		Help:    "Only print builtin resources"
		Long:    "builtin"
		Short:   ""
	}, {
		Name:    "custom"
		Type:    "bool"
		Default: "false"
		Help:    "Only print custom resources"
		Long:    "custom"
		Short:   ""
	}, {
		Name:    "local"
		Type:    "bool"
		Default: "false"
		Help:    "Only print workspace resources"
		Long:    "here"
		Short:   ""
	}]
}

#CreateCommand: schema.#Command & {
	TBD:   "α"
	Name:  "create"
	Usage: "create"
	Short: "create resources"
	Long:  Short
}

#GetCommand: schema.#Command & {
	TBD:   "α"
	Name:  "get"
	Usage: "get"
	Short: "find and display resources"
	Long:  Short
}

#EditCommand: schema.#Command & {
	TBD:   "α"
	Name:  "edit"
	Usage: "edit"
	Short: "edit resources"
	Long:  Short
}

#DeleteCommand: schema.#Command & {
	TBD:   "α"
	Name:  "delete"
	Usage: "delete"
	Short: "delete resources"
	Long:  Short
}
