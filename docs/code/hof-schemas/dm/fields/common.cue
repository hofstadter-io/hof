package fields

DataTypes: ID |
	UUID |
	CUID |
	Bool |
	String |
	Int |
	Float |
	Enum |
	Password |
	Email

ID: UUID & {Default: "" | *"uuid_generate_v4()"}

Field: {
	Name:   string
	Plural: string | *"\(Name)s"
	Type:   string
	Reln?:  string
}

UUID: Field & {
	Type:     "uuid"
	Nullable: bool | *false
	Unique:   bool | *true
	Default?: string
	Validation: {
		Format: "uuid"
	}
}

CUID: Field & {
	Type:     "cuid"
	Nullable: bool | *false
	Unique:   bool | *true
}

Bool: Field & {
	Type:     "bool"
	Default:  string | *"false"
	Nullable: bool | *false
}

String: Field & {
	Type:     "string"
	Length:   int | *64
	Unique:   bool | *false
	Nullable: bool | *false
	Default?: string
	Validation: {
		Max: int | *Length
	}
}

Int: Field & {
	Type:     "int"
	Nullable: bool | *false
	Default?: int
}

Float: Field & {
	Type:     "float"
	Nullable: bool | *false
	Default?: float
}

Enum: Field & {
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

Date: Field & {
	Type: "date"
}

Time: Field & {
	Type: "time"
}

Datetime: Field & {
	Type: "datetime"
}
