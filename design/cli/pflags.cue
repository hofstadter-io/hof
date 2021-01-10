package cli

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// TODO add: --non-intreactive (-y) ;  --no-color

#CliPflags: [...schema.#Flag] & [
		// Labels will be core
	//{
		//Name:    "labels"
		//Long:    "label"
		//Short:   "l"
		//Type:    "[]string"
		//Default: "nil"
		//Help:    "Labels for use across all commands"
	//},
	//{
		//Name:    "config"
		//Long:    "config"
		//Short:   "c"
		//Type:    "string"
		//Default: ""
		//Help:    "Path to a hof configuration file"
	//},
	//{
		//Name:    "secret"
		//Long:    "secret"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "The path to a hof secret file"
	//},
	//{
		//Name:    "contextFile"
		//Long:    "context-file"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "The path to a hof context file"
	//},
	//{
		//Name:    "context"
		//Long:    "context"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "The of an entry in the context file"
	//},
	//{
		//Name:    "global"
		//Long:    "global"
		//Short:   ""
		//Type:    "bool"
		//Default: "false"
		//Help:    "Operate using only the global config/secret context"
	//},
	//{
		//Name:    "local"
		//Long:    "local"
		//Short:   ""
		//Type:    "bool"
		//Default: "false"
		//Help:    "Operate using only the local config/secret context"
	//},

	// i/o formats and streams
	//{
		//Name:    "input"
		//Long:    "input"
		//Short:   "i"
		//Type:    "[]string"
		//Default: "nil"
		//Help:    "input streams, depending on the command context"
	//},
	//{
		//Name:    "inputFormat"
		//Long:    "input-format"
		//Short:   "I"
		//Type:    "string"
		//Default: ""
		//Help:    "input format, defaults to infered"
	//},

	//{
		//Name:    "output"
		//Long:    "output"
		//Short:   "o"
		//Type:    "[]string"
		//Default: "nil"
		//Help:    "output streams, depending on the command context"
	//},
	//{
		//Name:    "outputFormat"
		//Long:    "output-format"
		//Short:   "O"
		//Type:    "string"
		//Default: ""
		//Help:    "output format, defaults to cue"
	//},

	//{
		//Name:    "error"
		//Long:    "error"
		//Short:   ""
		//Type:    "[]string"
		//Default: "nil"
		//Help:    "error streams, depending on the command context"
	//},
	//{
		//Name:    "errorFormat"
		//Long:    "error-format"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "error format, defaults to cue"
	//},

	// context should encapsulate the next three
	//{
		//Name:    "account"
		//Long:    "account"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "the account context to use during this hof execution"
	//},
	//{
		//Name:    "billing"
		//Long:    "billing"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "the billing context to use during this hof execution"
	//},
	//{
		//Name:    "project"
		//Long:    "project"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "the project context to use during this hof execution"
	//},
	//{
		//Name:    "workspace"
		//Long:    "workspace"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "the workspace context to use during this hof execution"
	//},
	//{
		//Name:    "datamodelDir"
		//Long:    "datamodel-dir"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "directory for discovering resources"
	//},
	//{
		//Name:    "resourcesDir"
		//Long:    "resources-dir"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "directory for discovering resources"
	//},
	//{
		//Name:    "runtimesDir"
		//Long:    "runtimes-dir"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "directory for discovering runtimes"
	//},

	// these are more cue specific with a dash of hof
	{
		Name:    "package"
		Long:    "package"
		Short:   "p"
		Type:    "string"
		Default: ""
		Help:    "the Cue package context to use during execution"
	},
	//{
		//Name:    "errors"
		//Long:    "all-errors"
		//Short:   "E"
		//Type:    "bool"
		//Default: "false"
		//Help:    "print all available errors"
	//},
	//{
		//Name:    "ignore"
		//Long:    "ignore"
		//Short:   ""
		//Type:    "bool"
		//Default: "false"
		//Help:    "proceed in the presence of errors"
	//},
	//{
		//Name:    "simplify"
		//Long:    "simplify"
		//Short:   "S"
		//Type:    "bool"
		//Default: "false"
		//Help:    "simplify output"
	//},
	//{
		//Name:    "trace"
		//Long:    "trace"
		//Short:   ""
		//Type:    "bool"
		//Default: "false"
		//Help:    "trace cue computation"
	//},
	//{
		//Name:    "strict"
		//Long:    "strict"
		//Short:   ""
		//Type:    "bool"
		//Default: "false"
		//Help:    "report errors for lossy mappings"
	//},
	{
		Name:    "verbose"
		Long:    "verbose"
		Short:   "v"
		Type:    "int"
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
	//{
		//Name:    "ImpersonateAccount"
		//Long:    "impersonate-account"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "account to impersonate for this hof execution"
	//},
	//{
		//Name:    "traceToken"
		//Long:    "trace-token"
		//Short:   "T"
		//Type:    "string"
		//Default: ""
		//Help:    "used to help debug issues"
	//},
	//{
		//Name:    "LogHTTP"
		//Long:    "log-http"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "used to help debug issues"
	//},
	//{
		//Name:    "RunWeb"
		//Long:    "web"
		//Short:   ""
		//Type:    "bool"
		//Default: "false"
		//Help:    "run the command from the web ui"
	//},
	//{
		//Name:    "RunTUI"
		//Long:    "tui"
		//Short:   ""
		//Type:    "bool"
		//Default: "false"
		//Help:    "run the command from the terminal ui"
	//},
	//{
		//Name:    "RunREPL"
		//Long:    "repl"
		//Short:   ""
		//Type:    "bool"
		//Default: "false"
		//Help:    "run the command from the hof repl"
	//},
	//{
		//Name:    "Topic"
		//Long:    "topic"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "help topics for this command, 'list' will print available topics"
	//},
	//{
		//Name:    "Example"
		//Long:    "example"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "examples for this command, 'list' will print available examples"
	//},
	//{
		//Name:    "Tutorial"
		//Long:    "tutorial"
		//Short:   ""
		//Type:    "string"
		//Default: ""
		//Help:    "tutorials for this command, 'list' will print available tutorials"
	//},
]
