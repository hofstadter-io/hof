package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

EnvCommand: schema.Command & {
	Name:  "env"
	Usage: "env [args]"
	Short: "build, run, ship, and deploy environments (image, service, stack)"
	Long:  "build, run, ship, and deploy environments (image, service, stack)"

	// TODO, this should fallback to looking for commands in config, like cue & pnpm do
	// we want users to be able to define task/rules like Make, pnpm scripts, cue cmd
	// hof flow is this outside of the env/ci realm, we need something that mirrors onto dagger
	// this could also work into dagger modules/functions
	OmitRun: true

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
	}]

	// all subcommands get the name arg
	// Commands: [...{
	// 	Args: [{
	// 		Name:     "name"
	// 		Type:     "string"
	// 		Required: true
	// 		Help:     "name of the environment"
	// 	}, ..._]
	// }]

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
			Default: #"[]string{"local"}"# // todo, support special options like git-tag or git-commit
			Help:    "tags to give to the environment, can be set multiple times"
		}]
	}, {
		Name:  "info"
		Usage: "info [...target] [% ...cue]"
		Short: "get info about an environments"
		Long:  "get info about an environments"
	}, {
		Name:  "list"
		Usage: "list [...target] [% ...cue]"
		Short: "list environments"
		Long:  "list environments"
		Flags: [{
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
			Default: "nil"
			Help:    "sort columns, default is the order CUE defines"
		}]
	}, {
		Name:  "images"
		Usage: "images [...target] [% ...cue]"
		Short: "list images for an environments"
		Long:  "list images for an environments"
	}, {
		Name:  "ps"
		Usage: "ps [...target] [% ...cue]"
		Short: "print stats for running environments"
		Long:  "print stats for running environments"
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
		Name:  "down"
		Usage: "down [...target] [% ...cue]"
		Short: "stops an environment"
		Long:  "stops an environment"
	}, {
		Name:  "tag"
		Usage: "tag <src> <dst> [% ...cue]"
		Short: "tag an environment"
		Long:  "tag an environment"
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
	}, {
		Name:  "push"
		Usage: "push [...target] [% ...cue]"
		Short: "push an environment"
		Long:  "push an environment"
	}, {
		Name:  "pull"
		Usage: "pull [...target] [% ...cue]"
		Short: "pull an environment"
		Long:  "pull an environment"
	}, {
		Name:  "ci"
		Usage: "ci [...target] [% ...cue]"
		Short: "ci's an environment"
		Long:  "ci's an environment, local + remote parity"
	}, {
		Name:  "deploy"
		Usage: "deploy [...target] [% ...cue]"
		Short: "deploy an environment"
		Long:  "deploy an environment, think tf+helm"
	}]
}
