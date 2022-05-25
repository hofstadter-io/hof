package datamodel

import (
	"github.com/hofstadter-io/hof/schema/dm"
	"github.com/hofstadter-io/hof/schema/dm/fields"
	"github.com/hofstadter-io/hof/schema/dm/sql"
)

#MyModels: dm.#Datamodel & {
	Name: "MyModels"

	Models: {
		#AllModels
	}
}

#AllModels: dm.#Models & {
	"User":        User
	"UserProfile": UserProfile
}

User: dm.#Model & {
	Fields: {
		sql.#CommonFields
		email:   fields.Email
		persona: fields.Enum & {
			Vals: ["guest", "user", "admin", "owner"]
			Default: "guest"
		}
		password: fields.Password
		active:   fields.Bool
		username: fields.String
	}
}

UserProfile: dm.#Model & {
	Fields: {
		sql.#CommonFields
		firstName:  fields.String
		middleName: fields.String
		lastName:   fields.String
	}
}
