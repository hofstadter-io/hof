package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#ContextCommand: schema.#Command & {
	TBD:   "α"
	Name:  "context"
	Usage: "context"
	Short: "get, set, and use contexts"
	Long:  Short

	OmitRun: true

	Commands: [{
		TBD:   "β"
		Name:  "get"
		Usage: "get <key.path>"
		Short: "print a context or value(s) at path(s)"
		Long:  Short
	}, {
		TBD:   "β"
		Name:  "set"
		Usage: "set [expr]"
		Short: "set context values with an expr"
		Long:  Short
		Args: [{
			Name:     "expr"
			Type:     "string"
			Required: true
			Help:     "Cue expr for value you'd like to merge into your context"
		}]
	}, {
		TBD:   "Ø"
		Name:  "use"
		Usage: "use [name]"
		Short: "set a context as the current default"
		Long:  Short
	}, {
		TBD:   "Ø"
		Name:  "source"
		Usage: "source [name]"
		Short: "source a context into your environment"
		Long:  Short
	}, {
		TBD:   "Ø"
		Name:  "clear"
		Usage: "clear"
		Short: "clear your context and environment"
		Long:  Short
	}]
}

#ConfigCommand: schema.#Command & {
	// TBD:   "β"
	Name:  "config"
	Usage: "config"
	Short: "manage hof configuration"
	Long:  Short

	OmitRun: true

	Commands: [{
		// TBD:   "β"
		Name:  "get"
		Usage: "get <key.path>"
		Short: "print a config or value(s) at path(s)"
		Long:  Short
	}, {
		// TBD:   "β"
		Name:  "set"
		Usage: "set [expr]"
		Short: "set config values with an expr"
		Long:  Short
		Args: [{
			Name:     "expr"
			Type:     "string"
			Required: true
			Help:     "Cue expr for value you'd like to merge into your config"
		}]
	}, {
		// TBD:   "Ø"
		Name:  "use"
		Usage: "use [name]"
		Short: "bring a config into the current"
		Long:  Short
	}]
}

#SecretCommand: schema.#Command & {
	TBD:   "β"
	Name:  "secret"
	Usage: "secret"
	Short: "manage local secrets"
	Long:  Short

	OmitRun: true

	Commands: [{
		TBD:   "β"
		Name:  "get"
		Usage: "get <key.path>"
		Short: "print a secret or value(s) at path(s)"
		Long:  Short
	}, {
		TBD:   "β"
		Name:  "set"
		Usage: "set [expr]"
		Short: "set secret values with an expr"
		Long:  Short
		Args: [{
			Name:     "expr"
			Type:     "string"
			Required: true
			Help:     "Cue expr for value you'd like to merge into your secret"
		}]
	}, {
		TBD:   "Ø"
		Name:  "use"
		Usage: "use [name]"
		Short: "bring a secret into the current"
		Long:  Short
	}]
}
