package sql

import (
	"github.com/hofstadter-io/hof/schema/dm/fields"
)

#CommonFields: {
	ID:        fields.UUID
	CreatedAt: fields.Datetime
	UpdatedAt: fields.Datetime
	DeletedAt: fields.Datetime
}
