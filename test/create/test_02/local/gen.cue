package main

import (
	"github.com/hofstadter-io/hof/schema/common"
	"github.com/hofstadter-io/hof/schema/gen"
)

test: gen.#Generator & {
	@gen(test)

	Outdir: "./out"

	ModuleName: ""

	Create: {
		Message: {
			let name = Input.name
			Before: "testing creator before message"
			After:  "congrats, \(name) is ready, check \(Outdir)"
		}

		Input: _
		Prompt: [{
			Name:       "name"
			Type:       "input"
			Prompt:     "Please enter a name for..."
			Required:   true
			Validation: common.NameLabel
		}, {
			Name:       "about"
			Type:       "input"
			Prompt:     "Tell us about it..."
			Required:   true
			Validation: common.NameLabel
		}, {
			Name:   "frontend"
			Type:   "confirm"
			Prompt: "create frontend"
			Questions: [{
				Name:   "framework"
				Type:   "select"
				Prompt: "select framework"
				Options: ["React", "Vue", "Svelte"]
			}]
		}, {
			Name:   "backend"
			Type:   "confirm"
			Prompt: "create backend"
			Questions: [{
				Name:   "language"
				Type:   "select"
				Prompt: "select framework"
				Options: ["Go", "JS", "TS", "Python"]
			}]
		}, {
			Name:   "database"
			Type:   "confirm"
			Prompt: "create database"
			Questions: [{
				Name:   "vendor"
				Type:   "select"
				Prompt: "select framework"
				Options: ["Postgres", "Mysql", "Sqlite", "Mongo"]
			}]
		}, {
			Name:   "sdks"
			Type:   "confirm"
			Prompt: "create SDKs"
			Questions: [{
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

		PreFlow: {
			@flow()
			do: {
				@task(os.Stdout)
				text: "hello there!"
			}
		}
		PostFlow: {
			@flow()
			do: {
				@task(os.Exec)
				cmd: ["ls", "-lh"]
			}
		}
	}

	In: {
		Create.Input
		...
	}

	Out: [...gen.#File] & [ {
		TemplatePath: "debug"
		Filepath:     "debug.yaml"
	}, {
		TemplatePath: "readme.md"
		Filepath:     "readme.md"
	}, {
		TemplatePath: "gen.cue"
		Filepath:     "gen.cue"
	}]

	EmbeddedTemplates: {
		debug: {
			Content: """
				{{ yaml . }}
				"""
		}
	}
}
