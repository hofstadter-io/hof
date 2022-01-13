ver_20211230215755: {
	@datamodel(), @dm_ver(0.0.4)
	Name: "MyModels"

	// Models in the data model, ordered
	Models: {
		User: {
			@dm_model()
			Name: "User"
			Fields: {
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
				persona: {
					@dm_field()
					Name: "persona"
					Type: "string"
					Vals: ["guest", "user", "admin", "owner"]
					Nullable: false
					Default:  "guest"
				}
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
				ID: {
					@dm_field()
					Name:     "ID"
					Type:     "uuid"
					Nullable: false
					Unique:   true
					Default:  "gen_random_uuid()"
					Validation: {
						Format: "uuid"
					}
				}
				CID: {
					@dm_field()
					Name:     "CID"
					Type:     "cuid"
					Nullable: false
					Unique:   true
				}
				CreatedAt: {
					@dm_field()
					Name: "CreatedAt"
					Type: "datetime"
				}
				UpdatedAt: {
					@dm_field()
					Name: "UpdatedAt"
					Type: "datetime"
				}
				active: {
					@dm_field()
					Name:     "active"
					Type:     "bool"
					Default:  "false"
					Nullable: false
				}
				DeletedAt: {
					@dm_field()
					Name: "DeletedAt"
					Type: "datetime"
				}
			}
		}
		UserProfile: {
			@dm_model()
			Name: "UserProfile"
			Fields: {
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
				ID: {
					@dm_field()
					Name:     "ID"
					Type:     "uuid"
					Nullable: false
					Unique:   true
					Default:  "gen_random_uuid()"
					Validation: {
						Format: "uuid"
					}
				}
				CID: {
					@dm_field()
					Name:     "CID"
					Type:     "cuid"
					Nullable: false
					Unique:   true
				}
				CreatedAt: {
					@dm_field()
					Name: "CreatedAt"
					Type: "datetime"
				}
				UpdatedAt: {
					@dm_field()
					Name: "UpdatedAt"
					Type: "datetime"
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
				DeletedAt: {
					@dm_field()
					Name: "DeletedAt"
					Type: "datetime"
				}
			}
		}
	}
}
