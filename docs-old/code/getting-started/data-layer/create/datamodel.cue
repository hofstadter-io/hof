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
			}
		}
	}
}
