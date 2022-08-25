package main

import (
	"github.com/hofstadter-io/hof/schema/common"
	"github.com/hofstadter-io/hof/schema/gen"
)

#TestGen: gen.#HofGenerator & {
	@gen(test)

	Outdir: "./"

	PackageName: ""

	Create: [{
		Question:   "prompt"
		Name:       "name"
		Type:       "string"
		Default:    "./"
		Validation: common.NameLabel
	}]

	In: {
		foo: "bar"
		...
	}

	Out: [...gen.#HofGeneratorFile] & [
		{
			TemplatePath: "debug"
			Filepath:     "debug.yaml"
		},
	]

	EmbeddedTemplates: {
		debug: {
			Content: """
				{{ yaml . }}
				"""
		}
	}
}
