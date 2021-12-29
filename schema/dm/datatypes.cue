package dm

#DataTypes: #ID |
	#UUID |
	#CUID |
	#Bool |
	#String |
	#Enum |
	#Password |
	#Email

#ID: #UUID

#UUID: #Field & {
	Type: "uuid"
	Nullable: bool | *false
	Unique:   bool | *true
	Default: "gen_random_uuid()"
	Validation: {
		Format: "email"
	}
	...
}

#CUID: #Field & {
	Type: "cuid"
	Nullable: bool | *false
	Unique:   bool | *true
	...
}

#Bool: #Field & {
	Type: "bool"
	Default: string | *"false"
	Nullable: bool | *false
	...
}

#String: #Field & {
	Type: "string"
	Length: int | *64
	Unique: bool | *false
	Nullable: bool | *false
	Default?: string
	Validation: {
		Max: Length
	}
	...
}

#Enum: #Field & {
	Type: "string"
	Vals: [...string]
	Nullable: bool | *false
	Default?: string
	...
}

#Password: #String & {
	Bcrypt: true
}

#Email: #String & {
	Validation: {
		Format: "email"
	}
	Unique: true
	...
}

#Date: #Field & {
	Type: "date"
}

#Time: #Field & {
	Type: "time"
}

#Datetime: #Field & {
	Type: "datetime"
}
