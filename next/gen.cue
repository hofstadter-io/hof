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
		Name:   "docs"
		Module: "github.com/hofstadter-io/hof"
	}

	Datamodel: {}

}

Workflows: {
	dev: flows.dev
}
