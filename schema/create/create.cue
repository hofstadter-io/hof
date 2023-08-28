package create

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
	Prompt: [...Question]

	// (todo) Messages to print at start and end
	Message?: {
		Before: string
		After:  string
	}

	PreFlow:  _ // run hof flow beforehand
	PostFlow: _ // run hof flow afterwards
	// backwards compat
	PreFlow:  PreExec
	PreExec:  PreFlow
	PostFlow: PostExec
	PostExec: PostFlow

	// todo / consider
	// Check: _     // check for tools on host system
	// PreFlow: _   // run hof flow beforehand
	// PostFlow: _  // run hof flow afterwards
}

Question: {
	Name:   string
	Type:   "input" | "multiline" | "password" | "confirm" | "select" | "multiselect" | "subgroup"
	Prompt: string
	// for (multi)select
	Options?: [...string]
	Default?:    _
	Required:    bool | *false
	Validation?: _

	Questions?: [...Question]
}

// deprecated
#Creator:  Creator
#Question: Question
