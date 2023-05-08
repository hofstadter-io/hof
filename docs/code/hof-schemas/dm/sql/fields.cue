package sql

import (
	"github.com/hofstadter-io/hof/schema/dm/fields"
)

CommonFields: {
	ID:        fields.UUID & {Default: string | *"uuid_generate_v4()"}
	CreatedAt: fields.Datetime
	UpdatedAt: fields.Datetime
}

SoftDelete: {
	DeletedAt: fields.Datetime
}

PrimaryKey: fields.UUID & {
	Default: string | *"uuid_generate_v4()"
	SQL: PrimaryKey: true
}

Varchar: F=fields.String & {
	SQL: Type: "character varying(\(F.Length))"
}
