package example

// This is our input data, written as CUE
// The schema will be applied to validate
// the inputs and enrich the data model

Input: {
	User: {
		Fields: {
			id:       Type: "int"
			admin:    Type: "bool"
			username: Type: "string"
			email:    Type: "string"
		}
		Relations: {
			Profile: "HasOne"
			Post:    "HasMany"
		}
	}

	Profile: {
		Fields: {
			displayName: Type: "string"
			status:      Type: "string"
			about:       Type: "string"
		}
		Relations: User: "BelongsTo"
	}

	Post: {
		Fields: {
			title:  Type: "string"
			body:   Type: "string"
			public: Type: "bool"
		}
		Relations: User: "BelongsTo"
	}
}
