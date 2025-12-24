package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

EnvCommand: schema.Command & {
	Name:  "env"
	Usage: "env [args]"
	Short: "build, run, ship, and deploy environments (image, service, stack)"
	Long:  "build, run, ship, and deploy environments (image, service, stack)"

	OmitRun: true

	Pflags: [...schema.Flag] & [ {
		Name:    "Progress"
		Long:    "progress"
		Short:   "P"
		Type:    "string"
		Default: "\"auto\""
		Help:    "output format [auto, plain, tty, dots, report (for ai)]"
  },{
		Name:    "OnFailure"
		Long:    "on-failure"
		Short:   "F"
		Type:    "bool"
		Default: "false"
		Help:    "on failure, enter an interactive terminal, requires a tty"
  },{
		Name:    "NoExit"
		Long:    "no-exit"
		Short:   "T"
		Type:    "bool"
		Default: "false"
		Help:    "Leave the TUI open after finishing"
  },{
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
		Usage: "build [target...]"
		Short: "build an environment"
		Long: "build an environment"
	},{
		Name:  "export"
		Usage: "export [target...]"
		Short: "export an environment into local container runtime"
		Long: "export an environment into local container runtime"
	},{
		Name:  "info"
		Usage: "info [target...]"
		Short: "get info about an environments"
		Long: "get info about an environments"
	},{
		Name:  "list"
		Usage: "list"
		Short: "list environments"
		Long: "list environments"
    // flags for filtering
	},{
		Name:  "images"
		Usage: "images"
		Short: "list environments"
		Long: "list environments"
    // flags for filtering
	},{
		Name:  "ps"
		Usage: "ps [pattern...]"
		Short: "print stats for running environments"
		Long: "print stats for running environments"
    // flags for filtering
	},{
		Name:  "run"
		Usage: "run <name>"
		Short: "run an interactive environment"
		Long: "run an interactive environment"
		Args: [{
			Name:     "name"
			Type:     "string"
			Required: true
			Help:     "name of the environment"
		}]
	},{
		Name:  "up"
		Usage: "up [target...]"
		Short: "starts an environment"
		Long: "starts an environment"
	},{
		Name:  "down"
		Usage: "down [target...]"
		Short: "stops an environment"
		Long: "stops an environment"
	},{
		Name:  "tag"
		Usage: "tag <src> <dst>"
		Short: "tag an environment"
		Long: "tag an environment"
	},{
		Name:  "push"
		Usage: "push <name>"
		Short: "push an environment"
		Long: "push an environment"
	},{
		Name:  "pull"
		Usage: "pull <name>"
		Short: "pull an environment"
		Long: "pull an environment"
	},{
		Name:  "deploy"
		Usage: "deploy <name>"
		Short: "deploy an environment"
		Long: "deploy an environment"
	}]
}

