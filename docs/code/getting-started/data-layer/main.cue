package example

import (
	"github.com/hofstadter-io/hof/schema/dm/sql"
	"github.com/hofstadter-io/hof/schema/dm/fields"

	"hof.io/docs/example/gen"
)

// Generator definition
Migrations: gen.Migrations & {
	Outdir:      "out"
	"Datamodel": Datamodel
}

// The models we want in our database as tables & columns
Datamodel: sql.Datamodel & {
	$hof: metadata: name: "Datamodel"
	// these are the models for the application
	// they map onto database tables
	Models: {
		// Actual Models
		User: {
			Fields: {
				ID:        fields.UUID
				CreatedAt: fields.Datetime
				UpdatedAt: fields.Datetime
				DeletedAt: fields.Datetime

				email:    fields.Email & sql.Varchar
				password: fields.Password & sql.Varchar
				active:   fields.Bool
				verified: fields.Bool

				// this is the new field
				username: sql.Varchar

				// relation fields
				Profile: fields.UUID
				Profile: Relation: {
					Name:  "Profile"
					Type:  "belongs-to"
					Other: "Models.User"
				}
			}
		}

		UserProfile: {
			Fields: {
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
