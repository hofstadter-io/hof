package chat

import (
	"github.com/hofstadter-io/hof/schema"
)

// This is a complete Value tracked as one
// useful for schemas, config, and NoSQL
Chat: {
	schema.DHof// needed for reFerences
	$hof: chat: root: true

	Name:        string
	HumanName:   string | *Name
	MachineName: string | *Name

	Description:        string
	HumanDescription:   string | *Description
	MachineDescription: string | *Description

	Model: string

	// maps of named items for building up
	// input to an LLM
	Prompts: [string]:  string
	Examples: [string]: Example
	Messages: [string]: Message

	// map of named parameter sets
	Parameters: [string]: [string]: _
}

Example: {
	Input:  string
	Output: string
}

Message: {
	Role:    string
	Content: string
}
