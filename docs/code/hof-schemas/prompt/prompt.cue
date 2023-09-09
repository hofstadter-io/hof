package prompt

import (
	"github.com/hofstadter-io/hof/schema/common"
)

Prompt: {
	// the value / schema to be filled in
	// partially filled is supported, questions will be skipped
	Input: {...}

	// the prompt questions
	Questions: [...Question]

	// the resulting value after prompt
	Output: {...}
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

NamePrompt: {
	Name:       "name"
	Type:       "input"
	Prompt:     "What is your CLI named"
	Required:   true
	Validation: common.NameLabel
}

RepoPrompt: {
	Name:       "repo"
	Type:       "git"
	Prompt:     "Git repository"
	Validation: common.NameLabel
}
