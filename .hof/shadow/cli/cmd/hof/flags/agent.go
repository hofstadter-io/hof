package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var AgentFlagSet *pflag.FlagSet

type AgentFlagpole struct {
	Agent string
}

var AgentFlags AgentFlagpole

func SetupAgentFlags(fset *pflag.FlagSet, fpole *AgentFlagpole) {
	// flags

	fset.StringVarP(&(fpole.Agent), "agent", "A", "gemini-2.5-flash", "agent to use")
}

func init() {
	AgentFlagSet = pflag.NewFlagSet("Agent", pflag.ContinueOnError)

	SetupAgentFlags(AgentFlagSet, &AgentFlags)

}
