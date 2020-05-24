package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#ContextCommand: schema.#Command & {
	TBD:   "+ "
	Name:  "context"
	Usage: "context"
	Short: "Get, set, and use contexts"
	Long:  Short

	OmitRun: true

	Commands: [{
		Name:  "get"
		Usage: "get <key.path>"
		Short: "print a context or value(s) at path(s)"
		Long:  Short
	}, {
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
		TBD:   "+ "
		Name:  "use"
		Usage: "use [name]"
		Short: "set a context as the current default"
		Long:  Short
	}]
}

#ConfigCommand: schema.#Command & {
	Name:  "config"
	Usage: "config"
	Short: "Manage local configurations"
	Long:  Short

	OmitRun: true

	Commands: [{
		Name:  "get"
		Usage: "get <key.path>"
		Short: "print a config or value(s) at path(s)"
		Long:  Short
	}, {
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
	}]
}

#SecretCommand: schema.#Command & {
	Name:  "secret"
	Usage: "secret"
	Short: "Manage local secrets"
	Long:  Short

	OmitRun: true

	Commands: [{
		Name:  "get"
		Usage: "get <key.path>"
		Short: "print a secret or value(s) at path(s)"
		Long:  Short
	}, {
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
	}]
}
