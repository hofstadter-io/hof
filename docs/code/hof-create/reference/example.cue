package creator

import (
	"github.com/hofstadter-io/hof/schema/common"
	"github.com/hofstadter-io/hof/schema/gen"
)

Creator: gen.#Generator & {
	@gen(create)

	// the create configuration
	Create:

		// messages for the user
		Message: {
			// printed before the prompt
			Before: "A nice message before getting started"

			// printed after the prompt with Input available
			After:  """
			congrats, \(Input.name) is ready

			run the following command to get started
			
				hof flow
			"""
		}

		// provide a schema here
		// this is the value passed into the templates
		Input: _

		// the prompt, names and subq's should align with Input schema
		Prompt: [{

			// one line input prompt
			Name:       "name"
			Type:       "input"
			Prompt:     "Please enter a name for..."
			Required:   true
			Validation: common.NameLabel
		},{

			// Y/N confirmation
			Name:       "frontend"
			Type:       "confirm"
			Prompt:     "create frontend"

			// y/n can have subqustions
			Questions: [{
				// a single select prompt
				Name:   "framework"
				Type:   "select"
				Prompt: "select framework"
				Options: ["React", "Vue", "Svelt"]
			}]
		},{
			// another y/n prompt
			Name:       "sdks"
			Type:       "confirm"
			Prompt:     "create SDKs"
			Questions: [{
				// multi-select prompt
				Name:   "languages"
				Type:   "multiselect"
				Prompt: "select languages"
				Options: [
					"Go",
					"JavaScript",
					"Java",
					"Python",
					"Ruby",
					"Rust",
					"TypeScript",
				]
			}]
		}]
	}

	// data provided to the template system
	In: {
		// embed the Input results
		Create.Input
		...
	}

	// configure where to find templates, partials, and statics files
	gen.#SubdirTemplates & { #subdir: "creator" }

	// files to generate, relative to the <subdir>/templates/ directory
	Out: [...gen.#File] & [
		for file in [
			// javascript files
			"package.json",

			// cue module setup
			"cue.mods",
			"cue.mod/module.cue",

			// app config
			"app.cue",

			// a hof generator setup for sdks
			"sdks.cue"
		]{ TemplatePath: file, Filepath: file }
	]
}
