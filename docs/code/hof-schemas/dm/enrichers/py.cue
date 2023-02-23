package py

import "github.com/hofstadter-io/hof/schema/dm"

#FieldEnricher: {
	field: dm.#Field

	output: field
	output: PyType: [
		if field.Type == "uuid" { "uuid.UUID" },
		if field.Type == "date" { "datetime.date" },
		if field.Type == "time" { "datetime.time" },
		if field.Type == "datetime" { "datetime.datetime" },
		if field.Type == "string" { "str" },
		field.Type
	][0]
}
