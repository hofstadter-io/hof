package create


#Creator: {
	// schema and filled value for the create inputs
	// all inputs will be unified with this (files, flags, prompt)
	// it's contents should likely align with the prompt nesting
	CreateInput?: {...}

	// Init time inputs and prompts
	// if an entry will is already set by flags, it will be skipped
	CreatePrompt?: [...#Question]

	// (todo) Messages to print at start and end
	CreateMessage?: {
		Before: string
		After:  string
	}
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
