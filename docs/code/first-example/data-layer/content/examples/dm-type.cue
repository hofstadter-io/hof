package examples

import (
	"hof.io/docs/example/schema"
)

Datamodel: schema.Datamodel & {
	Name: "ExampleDatamodel"

	Models: {
		User: {
			Index: "Username"
			Fields: {
				Username: {Type: "string"}
				Email: {Type: "string"}
			}
		}

		Todo: {
			Index: "Title"
			Fields: {
				Title: {Type: "string"}
				Content: {Type: "string"}
			}
		}
	}
}
