ver_20211230215604: {
	@datamodel(), @dm_ver(0.0.2)
	Name: "MyModels"

	// Models in the data model, ordered
	Models: {
		User: {
			@dm_model()
			Name: "User"
			Fields: {
				// sql.#CommonFields
				email: {
					@dm_field()
					Name:     "email"
					Type:     "string"
					Length:   64
					Unique:   true
					Nullable: false
					Validation: {
						Max:    64
						Format: "email"
					}
				}
				//persona: fields.#Enum & {
				//vals: ["guest", "user", "admin", "owner"]
				//default: "guest"
				//} 
				password: {
					@dm_field()
					Name:     "password"
					Bcrypt:   true
					Type:     "string"
					Length:   64
					Unique:   false
					Nullable: false
					Validation: {
						Max: 64
					}
				}
				active: {
					@dm_field()
					Name:     "active"
					Type:     "bool"
					Default:  "false"
					Nullable: false
				}
			}
		}
		UserProfile: {
			@dm_model()
			Name: "UserProfile"
			Fields: {
				// sql.#CommonFields
				firstName: {
					@dm_field()
					Name:     "firstName"
					Type:     "string"
					Length:   64
					Unique:   false
					Nullable: false
					Validation: {
						Max: 64
					}
				}
				middleName: {
					@dm_field()
					Name:     "middleName"
					Type:     "string"
					Length:   64
					Unique:   false
					Nullable: false
					Validation: {
						Max: 64
					}
				}
				lastName: {
					@dm_field()
					Name:     "lastName"
					Type:     "string"
					Length:   64
					Unique:   false
					Nullable: false
					Validation: {
						Max: 64
					}
				}
			}
		}
	}
}
