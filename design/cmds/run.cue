package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

RunCommand: schema.Command & {
	Name:  "run"
	Usage: "run"
	Aliases: ["r"]
	Short: "Hof Line Script (HLS) is a successor to bash and python based scripting"
	Long:  RunCommandHelp

	Flags: [
		{
			Name:    "list"
			Type:    "bool"
			Default: "false"
			Help:    "list matching scripts that would run"
			Long:    "list"
			Short:   ""
			...
		},
		{
			Name:    "info"
			Type:    "bool"
			Default: "false"
			Help:    "view detailed info for matching scripts"
			Long:    "info"
			Short:   ""
			...
		},
		{
			Name:    "suite"
			Type:    "[]string"
			Default: "nil"
			Help:    "<name>: _ @run(suite)'s to run"
			Long:    "suite"
			...
		},
		{
			Name:    "runner"
			Type:    "[]string"
			Default: "nil"
			Help:    "<name>: _ @run(script)'s to run"
			Long:    "runner"
			...
		},
		{
			Name:    "environment"
			Type:    "[]string"
			Default: "nil"
			Help:    "exrta environment variables for scripts"
			Long:    "env"
			...
		},
		{
			Name:    "data"
			Type:    "[]string"
			Default: "nil"
			Help:    "exrta data to include in the scripts context"
			Long:    "data"
			...
		},
		{
			Name:    "workdir"
			Type:    "string"
			Default: ""
			Help:    "working directory"
			Long:    "workdir"
			...
		},
	]

}

RunCommandHelp: ##"""
	HofLineScript (HLS) run polyglot command and scripts seamlessly across runtimes
	
	can accept cue & flags or just a .hls file
	
	"""##
