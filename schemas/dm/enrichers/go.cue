package go

import "hofstadter.io/schemas/dm"

FieldEnricher: {
	field: dm.Field

	output: field
	output: GoType: [
			if field.Type == "uuid" {"uuid.UUID"},
			if field.Type == "datetime" {"time.Time"},
			if field.Type == "float" {"float64"},
			field.Type,
	][0]
}
