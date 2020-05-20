package cli

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#CliPflags: [...schema.#Flag] & [
		// Labels will be core
		{
		Name:    "labels"
		Long:    "label"
		Short:   "l"
		Type:    "[]string"
		Default: "nil"
		Help:    "Labels for use across all commands"
	},
	{
		Name:    "config"
		Long:    "config"
		Short:   ""
		Type:    "string"
		Default: ""
		Help:    "Path to a hof configuration file"
	},
	{
		Name:    "secret"
		Long:    "secret"
		Short:   ""
		Type:    "string"
		Default: ""
		Help:    "The path to a hof secret file"
	},
	{
		Name:    "context"
		Long:    "context"
		Short:   "C"
		Type:    "string"
		Default: ""
		Help:    "The path to a hof creds file"
	},
	{
		Name:    "global"
		Long:    "global"
		Short:   ""
		Type:    "bool"
		Default: "false"
		Help:    "Operate using only the global config/secret context"
	},
	{
		Name:    "local"
		Long:    "local"
		Short:   ""
		Type:    "bool"
		Default: "false"
		Help:    "Operate using only the local config/secret context"
	},

	// context should encapsulate the next three
	{
		Name:    "account"
		Long:    "account"
		Short:   ""
		Type:    "string"
		Default: ""
		Help:    "the account context to use during this hof execution"
	},
	{
		Name:    "billing"
		Long:    "billing"
		Short:   ""
		Type:    "string"
		Default: ""
		Help:    "the billing context to use during this hof execution"
	},
	{
		Name:    "project"
		Long:    "project"
		Short:   ""
		Type:    "string"
		Default: ""
		Help:    "the project context to use during this hof execution"
	},
	{
		Name:    "workspace"
		Long:    "workspace"
		Short:   ""
		Type:    "string"
		Default: ""
		Help:    "the workspace context to use during this hof execution"
	},
	// these are more cue specific with a dash of hof
	{
		Name:    "package"
		Long:    "package"
		Short:   "p"
		Type:    "string"
		Default: ""
		Help:    "the package context to use during this hof execution"
	},
	{
		Name:    "errors"
		Long:    "all-errors"
		Short:   "E"
		Type:    "bool"
		Default: "false"
		Help:    "print all available errors"
	},
	{
		Name:    "ignore"
		Long:    "ignore"
		Short:   ""
		Type:    "bool"
		Default: "false"
		Help:    "proceed in the presence of errors"
	},
	{
		Name:    "simplify"
		Long:    "simplify"
		Short:   "S"
		Type:    "bool"
		Default: "false"
		Help:    "simplify output"
	},
	{
		Name:    "trace"
		Long:    "trace"
		Short:   ""
		Type:    "bool"
		Default: "false"
		Help:    "trace computation"
	},
	{
		Name:    "strict"
		Long:    "strict"
		Short:   ""
		Type:    "bool"
		Default: "false"
		Help:    "report errors for lossy mappings"
	},
	{
		Name:    "verbose"
		Long:    "verbose"
		Short:   "v"
		Type:    "string"
		Default: ""
		Help:    "set the verbosity of output"
	},
	{
		Name:    "quiet"
		Long:    "quiet"
		Short:   "q"
		Type:    "bool"
		Default: ""
		Help:    "turn off output and assume defaults at prompts"
	},
	{
		Name:    "ImpersonateAccount"
		Long:    "impersonate-account"
		Short:   "I"
		Type:    "string"
		Default: ""
		Help:    "account to impersonate for this hof execution"
	},
	{
		Name:    "traceToken"
		Long:    "trace-token"
		Short:   "T"
		Type:    "string"
		Default: ""
		Help:    "used to help debug issues"
	},
	{
		Name:    "LogHTTP"
		Long:    "log-http"
		Short:   ""
		Type:    "string"
		Default: ""
		Help:    "used to help debug issues"
	},
	{
		Name:    "RunUI"
		Long:    "ui"
		Short:   ""
		Type:    "bool"
		Default: "false"
		Help:    "run the command from the web ui"
	},
	{
		Name:    "RunTUI"
		Long:    "tui"
		Short:   ""
		Type:    "bool"
		Default: "false"
		Help:    "run the command from the terminal ui"
	},
]
