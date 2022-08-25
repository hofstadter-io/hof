package gen

import (
  "github.com/hofstadter-io/hof/schema/common"
)

// Definition for a generator
#Generator: {
	// Base directory for the output
	Outdir: string | *"./"

	// Name of the generator, will default to kebab(label) where defined
	Name: common.NameLabel

	// Generater wide input value to templates.
	// Merged with any template or file level In values
	//   File.In will extend or replace any top-level fields here
	In: {...}

	// TODO, Generator wide cue.Value for writing incomplete values
	Val: {...}

	// File globs to watch and trigger regen when changed
	WatchFull: [...string]  // reloads & regens everything
	WatchFast:  [...string] // skips CUE reload, regens everything

	// Enable Diff3
	Diff3: bool | *true

	// Formatting Control
	Formatting: {
		// default for all files, unless overriden in a file
		Disabled: bool | *false

		// Should data files also be formatted?
		// (cue,yaml,json,toml,xml)
		FormatData: bool | *true

		// Map of names to formatter config values.
		//   Supports multiple configurations for a formatter,
		//   particularly useful for prettier.
		// Hof has defaults it will use if none are specified

		// map from file extensions to formatters
		Formatters: [Extenstion=string]: {
			// Name of the formatter, like 'prettier' or 'black'
			Formatter: string
			// formatter specific configuration
			Config: _
		}
	}

	// The final list of files for hof to generate
	Out: [...#File]

	// Template (top-level) TemplateConfig (globs+config)
	Templates: [...#Templates] | *[#Templates & {Globs: ["./templates/**/*"], TrimPrefix: "./templates/"}]

	// Partial (nested) TemplateConfig (globs+config)
	Partials: [...#Templates] | *[#Templates & {Globs: ["./partials/**/*"], TrimPrefix: "./partials/"}]

	// Statics are copied directly into the output, bypassing the rendering
	Statics: [...#Statics] | *[#Statics & {Globs: ["./static/**/*"], TrimPrefix: "./static/"}]

	// TODO, CUE files

	// The following mirror their non-embedded versions
	// however they have the content as a string in CUE
	// For templates and partials, Name is the path to reference
	EmbeddedTemplates: [name=string]: #Template
	EmbeddedPartials: [name=string]:  #Template
	// For statics, Name is the path to write the content
	EmbeddedStatics: [name=string]: string

	// TODO, consider adding 'Override*' for templates, partials, statics

	// For subgenerators so a generator can leverage and design for other hofmods
	Generators: [name=string]: #Generator & { Name: name }

	// schema and filled value for the create inputs
	// all inputs will be unified with this (files, flags, prompt)
	// it's contents should likely align with the prompt nesting
	CreateInput: {...}

	// Init time inputs and prompts
	// if an entry will is already set by flags, it will be skipped
	CreatePrompt: [...#CreateQuestion]

	// This should be set to default to the module name
	//   (i.e. 'string | *"github.com/<org>/<repo>"')
	// Users should not have to set this.
	// 
	// Used for indexing into the cue.mod/pkg directory...
	// until embed is supported, at which point this shouldn't be needed at all
	// only needed when you have example usage in the same module the generator is in
	// set to the empty string ("") as a generator writer who is making an example in the same module
	PackageName: string | *""
	// TODO, hof, can we introspect the generator / example packages and figure this out?

	// print debug info during load & gen
	Debug: bool | *false

	// Note, open so you can have any extra fields
	...
}

#CreateQuestion: {
	Name: string
	Type: "input" | "multiline" | "password" | "confirm" | "select" | "multiselect"
	Prompt: string
	// for (multi)select
	Options?: [...string]
	Default?: _
	Required: bool | *false
	Validation?: _

	Questions?: [...#CreateQuestion]
}

// deprecated
#HofGenerator: #Generator

