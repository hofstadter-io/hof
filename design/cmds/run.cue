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

	Flags: [ {
		Name:    "mode"
		Type:    "string"
		Default: "\"run\""
		Help:    "set the script execution mode"
		Long:    "mode"
		Short:   "m"
	}, {
		Name:    "workdir"
		Type:    "string"
		Default: ""
		Help:    "working directory"
		Long:    "workdir"
		Short:   "w"
	}, {
		Name:    "KeepTestdir"
		Type:    "bool"
		Default: "false"
		Help:    "keep the workdir after test mode run"
		Long:    "keep-testdir"
		Short:   ""
	}]
}

RunCommandHelp: ##"""
	Extended implementation of Testsuite	
	"""##
