package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

EnvCommand: schema.Command & {
	Name:  "env"
	Usage: "env [args]"
	Short: "build, run, ship, and deploy environments (image, service, stack)"
	Long:  "build, run, ship, and deploy environments (image, service, stack)"

	Pflags: [...schema.Flag] & [{
		Name:    "Progress"
		Long:    "progress"
		Short:   "P"
		Type:    "string"
		Default: "\"auto\""
		Help:    "output format [auto, plain, tty, dots, report (for ai)]"
	}, {
		Name:    "OnFailure"
		Long:    "on-failure"
		Short:   "F"
		Type:    "bool"
		Default: "false"
		Help:    "on failure, enter an interactive terminal, requires a tty"
	}, {
		Name:    "NoExit"
		Long:    "no-exit"
		Short:   "N"
		Type:    "bool"
		Default: "false"
		Help:    "Leave the TUI open after finishing"
	}, {
		Name:    "NoCache"
		Long:    "no-cache"
		Short:   "Z"
		Type:    "bool"
		Default: "false"
		Help:    "bust the cache and force evaluation"
	}, {
		Name:    "Shell"
		Long:    "shell"
		Short:   "S"
		Type:    "bool"
		Default: "false"
		Help:    "launch a shell after cmd completion for each target"
	}, {
		Name:    "Kind"
		Long:    "kind"
		Short:   "k"
		Type:    "[]string"
		Default: "nil"
		Help:    "kinds to include, defaults to all"
	}, {
		Name:    "Sort"
		Long:    "sort"
		Short:   "s"
		Type:    "[]string"
		Default: #"[]string{"name"}"# // todo, support special options like git-tag or git-commit
		Help:    "sort columns, default is the order CUE defines"
	}]

	Commands: [{
		Name:  "build"
		Usage: "build [...target] [% ...cue]"
		Short: "build an environment"
		Long:  "build an environment"
	}, {
		Name:  "export"
		Usage: "export [...target] [% ...cue]"
		Short: "export an environment into local container runtime"
		Long:  "export an environment into local container runtime"
		Flags: [{
			Name:    "Tag"
			Long:    "tag"
			Short:   "T"
			Type:    "[]string"
			Default: "nil" // todo, support special options like git-tag or git-commit
			Help:    "tags to give to the environment, can be set multiple times"
		}]
	}, {
		Name:  "get"
		Usage: "get [...target] [% ...cue]"
		Short: "get details for an environments"
		Long:  "get details for an environments"
	}, {
		Name:  "list"
		Usage: "list [...target] [% ...cue]"
		Short: "list environments"
		Long:  "list environments"
	}, {
		Name:  "run"
		Usage: "run <target> [% [...cue]]"
		Short: "run an interactive environment"
		Long:  "run an interactive environment"
		Flags: [{
			Name:    "Command"
			Long:    "cmd"
			Short:   "c"
			Type:    "string"
			Default: "\"\""
			Help:    "the command to run"
		}]
	}, {
		Name:  "up"
		Usage: "up [...target] [% ...cue]"
		Short: "starts an environment"
		Long:  "starts an environment"
	}, {
		Name:  "publish"
		Usage: "publish [...target] [% ...cue]"
		Short: "publish an environment"
		Long:  "publish an environment"
		Flags: [{
			Name:    "Registry"
			Long:    "registry"
			Short:   "R"
			Type:    "string"
			Default: #""host.docker.internal:5000""#
			Help:    "registry to push to, defaults to veg internal"
		}, {
			Name:    "Tag"
			Long:    "tag"
			Short:   "T"
			Type:    "[]string"
			Default: #"[]string{"local"}"# // todo, support special options like git-tag or git-commit
			Help:    "tags to give to the environment, can be set multiple times"
		}]
	}]
}
