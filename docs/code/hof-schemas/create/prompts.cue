package create

import (
	"github.com/hofstadter-io/hof/schema/common"
)

#NamePrompt: {
	Name:       "name"
	Type:       "input"
	Prompt:     "What is your CLI named"
	Required:   true
	Validation: common.NameLabel
}

#RepoPrompt: {
	Name:       "repo"
	Type:       "git"
	Prompt:     "Git repository"
	Validation: common.NameLabel
}
