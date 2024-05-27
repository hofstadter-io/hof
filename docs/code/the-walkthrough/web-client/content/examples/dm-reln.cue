ServerDatamodel: schema.#Datamodel & {
	Models: {
		User: {
			Relations: {
				Todos: {
					Reln: "HasMany"
					Type: "Todo"
				}
			}
		}
		Todo: {
			Relations: {
				Author: {
					Reln: "BelongsTo"
					Type: "User"
				}
			}
		}
	}
}
