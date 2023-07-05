package docs

import (
	"github.com/hofstadter-io/supacode/flows"
	"github.com/hofstadter-io/supacode/gen"
	"github.com/hofstadter-io/supacode/schema"
)

// Generator definition
Generator: gen.Generator & {
	Name:   "docs"
	Outdir: "./"
	App:    schema.App & {
		name:   "docs"
		module: "github.com/hofstadter-io/hof"

		search: enabled:   false
		auth: enabled:     false
		database: enabled: false
	}

	Datamodel: {}

}

Workflows: {
	dev: flows.dev
}
