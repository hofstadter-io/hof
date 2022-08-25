package main

import (
	"github.com/hofstadter-io/hof/schema/common"
	"github.com/hofstadter-io/hof/schema/gen"
)

test: gen.#HofGenerator & {
	@gen(test)

	Outdir: "./"

	PackageName: ""

	CreateInput: {
		name: string
		about: string
		frontend?: {
			framework: string
		}
		backend?: {
			language: string
		}
		database?: {
			vendor: string
		}
		sdks?: {
			languages: [...string]
		}
	}

	CreatePrompt: [{
		Name:       "name"
		Type:       "input"
		Prompt:     "Please enter a name for..."
		Required:   true
		Validation: common.NameLabel
	},{
		Name:       "about"
		Type:       "multiline"
		Prompt:     "Tell us about it..."
		Required:   true
		Validation: common.NameLabel
	},{
		Name:       "frontend"
		Type:       "confirm"
		Prompt:     "create frontend"
		Questions: [{
			Name:   "framework"
			Type:   "select"
			Prompt: "select framework"
			Options: ["React", "Vue", "Svelt"]
		}]
	},{
		Name:       "backend"
		Type:       "confirm"
		Prompt:     "create backend"
		Questions: [{
			Name:   "language"
			Type:   "select"
			Prompt: "select framework"
			Options: ["Go", "JS", "TS", "Python"]
		}]
	},{
		Name:       "database"
		Type:       "confirm"
		Prompt:     "create database"
		Questions: [{
			Name:   "vendor"
			Type:   "select"
			Prompt: "select framework"
			Options: ["Postgres", "Mysql", "Sqlite", "Mongo"]
		}]
	},{
		Name:       "sdks"
		Type:       "confirm"
		Prompt:     "create SDKs"
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

	In: {
		CreateInput
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
