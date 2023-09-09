package create

import (
	"github.com/hofstadter-io/hof/schema/prompt"
)

Creator: Create: {
	// schema and filled value for the create inputs
	// all inputs will be unified with this (files, flags, prompt)
	// it's contents should likely align with the prompt nesting
	Input: {...}

	// extra args provided by user when calling `hof create <repo> [args]`
	// Filled by hof, you can map these to inputs however you like
	Args: [...string]

	// Init time inputs and prompts
	// if an entry will is already set by flags, it will be skipped
	Questions: [...prompt.Question]
	Questions: Prompt
	Prompt:    Questions

	// (todo) Messages to print at start and end
	Message?: {
		Before: string
		After:  string
	}

	PreFlow?:  _ // run hof flow beforehand
	PostFlow?: _ // run hof flow afterwards

	// todo / consider
	// Check: _     // check for tools on host system
	// PreFlow: _   // run hof flow beforehand
	// PostFlow: _  // run hof flow afterwards
}

// deprecated
#Creator:   Creator
#Question:  Question
Question:   prompt.Question
NamePrompt: prompt.NamePrompt
RepoPrompt: prompt.RepoPrompt
