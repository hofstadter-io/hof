package examples

import (
	"hof.io/docs/example/schema"
)

Datamodel: schema.Datamodel & {
	@datamodel(datamodel)
	$hof: metadata: name: "Datamodel"

	Models: {
		User: {
			Index: "Username"
			Fields: {
				Username: {Type: "string"}
				Email: {Type: "string"}
			}
			Relations: {
				Todos: {
					Reln: "HasMany"
					Type: "Todo"
				}
			}
		}

		Todo: {
			Index: "Title"
			Fields: {
				Title: {Type: "string"}
				Content: {Type: "string"}
			}
			Relations: {
				Author: {
					Reln: "BelongsTo"
					Type: "User"
				}
			}
		}
	}
}
