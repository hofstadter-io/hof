package examples

import (
	"hof.io/docs/example/schema"
)

Datamodel: schema.Datamodel & {
	@datamodel(datamodel)
	#hof: metadata: name: "Datamodel"

	Models: {
		User: {
			Index: "Username"
			Fields: {
				Username: {Type: "string"}
				Email: {Type: "string"}
				Todos: {
					Type: "string"
					Relation: {
						Type:  "has-many"
						Name:  "Todos"
						Other: "Todo"
					}
				}
			}
		}

		Todo: {
			Index: "Title"
			Fields: {
				Title: {Type: "string"}
				Content: {Type: "string"}
				Author: {
					Type: "string"
					Relation: {
						Type:  "belongs-to"
						Name:  "Author"
						Other: "User"
					}
				}
			}
		}
	}
}
