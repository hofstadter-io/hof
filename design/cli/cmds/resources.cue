package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// Kubernetes inspired commands (maybe some hyper-cloud too)

#InfoCommand: schema.#Command & {
	TBD:   "α"
	Name:  "info"
	Usage: "info"
	Aliases: ["i"]
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
	Aliases: ["c"]
	Short: "create resources"
	Long:  Short
}

#GetCommand: schema.#Command & {
	TBD:   "α"
	Name:  "get"
	Usage: "get"
	Aliases: ["g"]
	Short: "find and display resources"
	Long:  Short
}

#SetCommand: schema.#Command & {
	TBD:   "α"
	Name:  "set"
	Usage: "set"
	Aliases: ["s"]
	Short: "find and configure resources"
	Long:  Short
}

#EditCommand: schema.#Command & {
	TBD:   "α"
	Name:  "edit"
	Usage: "edit"
	Aliases: ["e"]
	Short: "edit resources"
	Long:  Short
}

#DeleteCommand: schema.#Command & {
	TBD:   "α"
	Name:  "delete"
	Usage: "delete"
	Aliases: ["del"]
	Short: "delete resources"
	Long:  Short
}
