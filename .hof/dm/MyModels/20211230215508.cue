ver_20211230215508: {
	@datamodel(), @dm_ver(0.0.1)
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
		// "UserProfile": UserProfile
	}
}
