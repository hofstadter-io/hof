package agent

import (
	"github.com/hofstadter-io/hof/schemas"
	"github.com/hofstadter-io/hof/schemas/common"
)

// Chat represents a call to an LLM via the `hof chat with' command.
// You can put these in your module to provide ChatGPT like interactions
// for the other components in your module, or make a module just for Chats.
Agent: {
	schema.Hof // needed for reFerences

	#hof: agent: root: true

	Name: common.NameLabel
}
