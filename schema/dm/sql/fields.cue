package sql

import (
	"github.com/hofstadter-io/hof/schema/dm/fields"
)

CommonFields: {
	ID:        fields.UUID & {Default: "" | *"uuid_generate_v4()"}
	CreatedAt: fields.Datetime
	UpdatedAt: fields.Datetime
}

SoftDelete: {
	DeletedAt: fields.Datetime
}

Varchar: F=fields.String & {
	sqlType: "varchar(\(F.Length))"
}
