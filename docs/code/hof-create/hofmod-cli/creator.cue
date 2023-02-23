package creators

import (
	"github.com/hofstadter-io/hof/schema/common"
	"github.com/hofstadter-io/hof/schema/gen"
)

// A creator is a generator
Creator: gen.#Generator & {
	@gen(creator)

	// The create section
	Create: {

		// pre|post messages
		Message: {
			// printed before input prompting
			Before: "Creating a new Go Cli"

			let name = Input.name
			// printed after successfully running
			After: """
			Your new Cli generator is ready, run the following
			to generate the code, build the binary, and run \(name).

			now run 'make first'    (cd to the --outdir if used)
			"""
		}

		// Input schema
		Input: {
			name:      string
			repo:      string
			about:     string
			releases:  bool | *false
			updates:   bool | *false
			telemetry: bool | *false
		}

		// Prompt configuration, a list of questions
		Prompt: [{
			Name:       "name"
			Type:       "input"
			Prompt:     "What is your CLI named"
			Required:   true
			Validation: common.NameLabel
		},{
			Name:       "repo"
			Type:       "input"
			Prompt:     "Git repository"
			Default:    "github.com/user/repo"
			Validation: common.NameLabel
		},{
			Name:       "about"
			Type:       "input"
			Prompt:     "Tell us a bit about it..."
			Required:   true
			Validation: common.NameLabel
		},{
			Name:       "releases"
			Type:       "confirm"
			Prompt:     "Enable GoReleaser tooling"
			Default:    true
		},

		// conditional prompts based on prior inputs
		if Input.releases == true {
			Name:       "updates"
			Type:       "confirm"
			Prompt:     "Enable self updating"
			Default:    true
		}

		if Input.releases == true {
			Name:       "telemetry"
			Type:       "confirm"
			Prompt:     "Enable telemetry"
		}
		]
	}

	// The user input is embedded the only input
	In: {
		Create.Input
	}

	// the commont template directories can be found under
	// the directory '{repo}/creators/{templates,partials,...}'
	gen.#SubdirTemplates & { #subdir: "creators" }

	// the files which will be generated
	Out: [...gen.#File] & [ 
		for file in [
			"cli.cue",            // starting CUE for the cli code generator
			"cue.mods",           // CUE module file
			"cue.mod/module.cue", // CUE module file
			"Makefile",           // for an easy post-create command for the user
		]{ TemplatePath: file, Filepath: file }
	]
}
