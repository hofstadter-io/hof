package create


#Creator: Create: {
	// schema and filled value for the create inputs
	// all inputs will be unified with this (files, flags, prompt)
	// it's contents should likely align with the prompt nesting
	Input: {...}

	// Init time inputs and prompts
	// if an entry will is already set by flags, it will be skipped
	Prompt: [...#Question]

	// (todo) Messages to print at start and end
	Message?: {
		Before: string
		After:  string
	}

	// todo / consider
	// Check: _     // check for tools on host system
	// PreFlow: _   // run hof flow beforehand
	// PostFlow: _  // run hof flow afterwards
}

#Question: {
	Name: string
	Type: "input" | "multiline" | "password" | "confirm" | "select" | "multiselect"
	Prompt: string
	// for (multi)select
	Options?: [...string]
	Default?: _
	Required: bool | *false
	Validation?: _

	Questions?: [...#Question]
}
