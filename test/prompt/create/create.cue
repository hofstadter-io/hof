package create

import (
	"github.com/hofstadter-io/hof/schema/common"
	"github.com/hofstadter-io/hof/schema/gen"
)

Creator: gen.Generator & {
	@gen(creator)

	Create: {
		Message: {
			let name = Input.name
			Before: "This is a test prompt"
			After:  """
			Created project: \(name).
			"""
		}

		Input: {
			name:     string
			repo:     string
			releases: bool | *false
		}

		// Questions: Prompt
		Prompt: [{
			Name:       "name"
			Type:       "input"
			Prompt:     "What is your project named"
			Required:   true
			Validation: common.NameLabel
		}, {
			Name:       "repo"
			Type:       "input"
			Prompt:     "Git repository"
			Default:    "github.com/user/repo"
			Validation: common.NameLabel
		}, {
			Name:    "releases"
			Type:    "confirm"
			Prompt:  "Enable release tooling"
			Default: true
		}]
	}

	In: {
		Create.Input
		...
	}

	Out: [...gen.File] & [{
		Filepath:     "debug.yaml"
		TemplatePath: "debug.yaml"
	}]

	Statics: []
	Partials: []
	Templates: []

	EmbeddedTemplates: {
		"debug.yaml": {
			Content: """
				{{ yaml . }}
				"""
		}
	}

	ModuleName: ""
}
