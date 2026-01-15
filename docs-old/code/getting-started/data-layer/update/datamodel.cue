package datamodel

import (
	"github.com/hofstadter-io/hof/schema/dm/sql"
	"github.com/hofstadter-io/hof/schema/dm/fields"
)

// Traditional database model which maps onto tables & columns.
Datamodel: sql.Datamodel & {
	Models: {
		User: {
			Fields: {
				ID:        fields.UUID
				CreatedAt: fields.Datetime
				UpdatedAt: fields.Datetime
				DeletedAt: fields.Datetime

				email:    fields.Email
				username: fields.String
				password: fields.Password
				verified: fields.Bool
				active:   fields.Bool

				persona: fields.Enum & {
					Vals: ["guest", "user", "admin", "owner"]
					Default: "user"
				}

				// relation fields
				Profile: fields.UUID
				Profile: Relation: {
					Name:  "Profile"
					Type:  "has-one"
					Other: "Models.UserProfile"
				}
			}
		}

		UserProfile: {
			Fields: {
				// note how we are using sql fields here
				sql.CommonFields

				About:  sql.Varchar
				Avatar: sql.Varchar
				Social: sql.Varchar

				Owner: fields.UUID
				Owner: Relation: {
					Name:  "Owner"
					Type:  "belongs-to"
					Other: "Models.User"
				}
			}
		}
	}
}
