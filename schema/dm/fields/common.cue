package fields

import (
	"github.com/hofstadter-io/hof/schema/dm"
)

#DataTypes: ID |
	UUID |
	CUID |
	Bool |
	String |
	Int |
	Float |
	Enum |
	Password |
	Email

ID: UUID

UUID: dm.#Field & {
	Type:     "uuid"
	Nullable: bool | *false
	Unique:   bool | *true
	Default:  string | *"gen_random_uuid()"
	Validation: {
		Format: "uuid"
	}
}

CUID: dm.#Field & {
	Type:     "cuid"
	Nullable: bool | *false
	Unique:   bool | *true
}

Bool: dm.#Field & {
	Type:     "bool"
	Default:  string | *"false"
	Nullable: bool | *false
}

String: dm.#Field & {
	Type:     "string"
	Length:   int | *64
	Unique:   bool | *false
	Nullable: bool | *false
	Default?: string
	Validation: {
		Max: Length
	}
}

Int: dm.#Field & {
	Type: "int"
	Nullable: bool | *false
	Default?: int
}

Float: dm.#Field & {
	Type: "float"
	Nullable: bool | *false
	Default?: float
}

Enum: dm.#Field & {
	Type: "string"
	Vals: [...string]
	Nullable: bool | *false
	Default?: string
}

Password: String & {
	Bcrypt: true
}

Email: String & {
	Validation: {
		Format: "email"
	}
	Unique: true
}

Date: dm.#Field & {
	Type: "date"
}

Time: dm.#Field & {
	Type: "time"
}

Datetime: dm.#Field & {
	Type: "datetime"
}
