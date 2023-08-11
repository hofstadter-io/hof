package schema

import (
	"github.com/hofstadter-io/hof/schema/dm/sql"
	"github.com/hofstadter-io/hof/schema/dm/enrichers:go"
	"github.com/hofstadter-io/hof/schema/dm/enrichers:py"
)

Datamodel: sql.Datamodel & {
	Models: {
		@history()

		// apply to each "model" (CUE pattern constraint)
		[string]: {
			Fields: {
				@history()

				// apply to each "field" (CUE pattern constraint)
				[string]: go.FieldEnricher
				[string]: py.FieldEnricher

				// These will add GoType and PyType to the "model" "field"
				...
			}
		}
		...
	}
}
