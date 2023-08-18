package main

import (
	"github.com/hofstadter-io/hof/schema/dm/sql"
	"github.com/hofstadter-io/supacode/gen"
	"github.com/hofstadter-io/supacode/schema"
)

MyGen: gen.Generator & {
	@gen()
	Datamodel: MyModels
	App:       MyApp
}

MyModels: sql.Datamodel & {
	@datamodel()
	Models: {
		...
	}
}

MyApp: schema.App & {
	// sometimes a module author will place the datamodel here
	// or provide a custom datamodel specific to their generator
	...
}
