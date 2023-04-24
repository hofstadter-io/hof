package gen

import (
	"github.com/hofstadter-io/hof/schema/dm/sql"
	"github.com/hofstadter-io/hof/schema/gen"
)

// Generator definition
Migrations: gen.#Generator & {

	// User inputs to this generator
	// -----------------------------

	// The server design conforming to the server schema
	Datamodel: sql.Datamodel

	// Base output directory, defaults to current
	Outdir: string | *"./"

	// Required fields for hof
	// ------------------------

	// In is passed to every template
	In: {
		DM: Datamodel

		// this will cause an issue during injection
		// we won't know what datamodel this local 'In.User' really belongs to
		// this applies at the File.In scope as well
		// User: Datamodel.Models.User
	}

	Statics: []

	// Actual files generated by hof, combined into a single list
	Out: [...gen.#File] & _All

	_All: [
		for _, F in _OnceFiles {F},
		for _, F in _MigrationFiles {F},
	]

	// Note, we can omit Templates, Partials, and Statics
	// since the default values are sufficient for us

	// Internal fields for mapping Input to templates
	// ----------------------------------------------

	// Files that are generated once per server
	_OnceFiles: [...gen.#File] & [
			{
			TemplatePath: "Makefile"
			Filepath:     "Makefile"
		},
		{
			TemplatePath: "debug.txt"
			Filepath:     "debug.txt"
		},
		{
			TemplatePath: "migration.sql"
			Filepath:     "migrations/latest.sql"
		},
	]

	_MigrationFiles: [...gen.#File] & [ for _, S in Datamodel.History {
		In: {
			Snapshot: S
		}
		TemplatePath: "migration.sql"
		Filepath:     "migrations/\(S.Timestamp).sql"
	}]
	...
}
