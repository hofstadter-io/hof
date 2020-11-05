package schema

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
	type: "uuid"
	nullable: bool | *false
	unique:   bool | *true
	generate: bool | *true
	...
}

#CUID: #Field & {
	type: "cuid"
	nullable: bool | *false
	unique:   bool | *true
	generate: bool | *true
	...
}

#Bool: #Field & {
	type: "bool"
	default: bool | *false
	nullable: bool | *false
	...
}

#String: #Field & {
	type: "string"
	length: int | *64
	unique: bool | *false
	nullable: bool | *false
	default?: string
	...
}

#Enum: #Field & {
	type: "string"
	vals: [...string]
	nullable: bool | *false
	default?: string
	...
}

#Password: #Field & {
	type: "text"
	bcrypt: true
}

#Email: #String & {
	validation: "email"
	...
}

#Date: #Field & {
	type: "date"
}

#Time: #Field & {
	type: "time"
}

#Datetime: #Field & {
	type: "datetime"
}
