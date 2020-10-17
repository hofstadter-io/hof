package schema

#DataTypes: #ID |
	#UUID |
	#CUID |
	#XUID |
	#Bool |
	#String |
	#Enum |
	#Password |
	#Email

#ID: #UUID

#UUID: {
	type: "uuid"
	nullable: bool | *false
	unique:   bool | *true
	generate: bool | *true
	...
}

#CUID: {
	type: "cuid"
	nullable: bool | *false
	unique:   bool | *true
	generate: bool | *true
	...
}

#XUID: {
	type: "xuid"
	nullable: bool | *false
	unique:   bool | *true
	generate: bool | *true
	...
}

#Bool: {
	type: "boolean"
	default: bool | *false
	nullable: bool | *false
	...
}

#String: {
	type: "string"
	length: int | *64
	unique: bool | *false
	nullable: bool | *false
	default?: string
	...
}

#Enum: {
	type: "string"
	vals: [...string]
	nullable: bool | *false
	default?: string
	...
}

#Password: {
	type: "text"
	bcrypt: true
}

#Email: #String & {
	validation: "email"
	...
}

