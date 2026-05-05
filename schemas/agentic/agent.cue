package agentic

import (
	"github.com/hofstadter-io/hof/schemas"
	"github.com/hofstadter-io/hof/schemas/common"
)

// Chat represents a call to an LLM via the `hof chat with' command.
// You can put these in your module to provide ChatGPT like interactions
// for the other components in your module, or make a module just for Chats.
Agent: {
	schema.Hof // needed for reFerences

  #hof: agentic: { root: true, kind: "agent" }

	name: common.NameLabel

	// default model
	model?: string

	// description for man and machine
	description?: string

	// content or file path
	instruction: string

	tools: [...string]
	toolsets: [...string]
	mcp: [...string]
	subagents: [...string]

	// name of an agentic environment
	environ?: string
	
}
