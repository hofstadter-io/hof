package cli

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#CliPflags: [...schema.#Flag] & [
  {
    Name:    "config"
    Long:    "config"
    Short:   "C"
    Type:    "string"
    Default: ""
    Help:    "Path to a hof configuration file"
  },
	{
    Name:    "context"
    Long:    "context"
    Short:   "X"
    Type:    "string"
    Default: ""
    Help:    "The path to a hof creds file"
  },

	// context should encapsulate the next three
  {
    Name:    "account"
    Long:    "account"
    Short:   "A"
    Type:    "string"
    Default: ""
    Help:    "the account context to use during this hof execution"
  },
  {
    Name:    "billing"
    Long:    "billing"
    Short:   "B"
    Type:    "string"
    Default: ""
    Help:    "the billing context to use during this hof execution"
  },
  {
    Name:    "project"
    Long:    "project"
    Short:   "P"
    Type:    "string"
    Default: ""
    Help:    "the project context to use during this hof execution"
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
		Name:    "labels"
		Long:    "label"
		Short:   "L"
		Type:    "[]string"
		Default: "nil"
		Help:    "Labels for use across all commands"
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
		Short:   "i"
		Type:    "bool"
		Default: "false"
		Help:    "proceed in the presence of errors"
	},
	{
		Name:    "simplify"
		Long:    "simplify"
		Short:   "s"
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
    Short:   ""
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

	// Keep this around for DMA legacy for the moment
	{
    Name:    "creds"
    Long:    "creds"
    Short:   ""
    Type:    "string"
    Default: ""
    Help:    "The path to a hof creds file"
  },
]

