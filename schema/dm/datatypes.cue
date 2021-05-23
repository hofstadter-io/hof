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
	default: "gen_random_uuid()"
	validation: {
		format: "email"
	}
	...
}

#CUID: #Field & {
	type: "cuid"
	nullable: bool | *false
	unique:   bool | *true
	...
}

#Bool: #Field & {
	type: "bool"
	default: string | *"false"
	nullable: bool | *false
	...
}

#String: #Field & {
	type: "string"
	length: int | *64
	unique: bool | *false
	nullable: bool | *false
	default?: string
	validation: {
		max: length
	}
	...
}

#Enum: #Field & {
	type: "string"
	vals: [...string]
	nullable: bool | *false
	default?: string
	...
}

#Password: #String & {
	bcrypt: true
}

#Email: #String & {
	validation: {
		format: "email"
	}
	unique: true
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
