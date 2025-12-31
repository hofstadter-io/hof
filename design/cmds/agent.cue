package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

AgentCommand: schema.Command & {
	Name:  "agent"
	Usage: "agent [args]"
	Short: "run an agent"
	Long:  AgentRootHelp

	Flags: [{
		Name:    "agent"
		Type:    "string"
		Default: "\"gemini-2.5-flash\""
		Help:    "agent to use"
		Long:    "agent"
		Short:   "A"
	}]

}

AgentRootHelp: #"""
	Run an agent
	"""#
