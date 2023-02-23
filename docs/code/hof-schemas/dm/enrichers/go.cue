package go

import "github.com/hofstadter-io/hof/schema/dm"

#FieldEnricher: {
	field: dm.#Field

	output: field
	output: GoType: [
		if field.Type == "uuid" { "uuid.UUID" },
		if field.Type == "datetime" { "time.Time" },
		if field.Type == "float" { "float64" },
		field.Type,
	][0]
}
