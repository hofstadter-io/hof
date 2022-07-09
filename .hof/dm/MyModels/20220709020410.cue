ver_20220709020410: {
	@hof(datamodel), @dm_ver(0.0.5)
	Name: "MyModels"

	// Models in the data model, ordered
	Models: {
		User: {
			@hof(model)
			Fields: {
				email: {
					@hof(field)

					// this should be a string you can use within your templates
					Type:   "string"
					Name:   "email"
					Length: 64
					Labels: {}
					Unique:   true
					Nullable: false
					Validation: {
						Max:    64
						Format: "email"
					}
				}
				persona: {
					@hof(field)

					// this should be a string you can use within your templates
					Type: "string"
					Name: "persona"
					Vals: ["guest", "user", "admin", "owner"]
					Nullable: false
					Labels: {}
					Default: "guest"
				}
				password: {
					@hof(field)

					// this should be a string you can use within your templates
					Type:     "string"
					Name:     "password"
					Length:   64
					Unique:   false
					Nullable: false
					Bcrypt:   true
					Labels: {}
					Validation: {
						Max: 64
					}
				}
				ID: {
					@hof(field)

					// this should be a string you can use within your templates
					Type:     "uuid"
					Name:     "ID"
					Nullable: false
					Unique:   true
					Default:  "gen_random_uuid()"
					Labels: {}
					Validation: {
						Format: "uuid"
					}
				}
				CID: {
					@hof(field)

					// this should be a string you can use within your templates
					Type:     "cuid"
					Name:     "CID"
					Nullable: false
					Labels: {}
					Unique: true
				}
				CreatedAt: {
					@hof(field)

					// this should be a string you can use within your templates
					Type: "datetime"
					Name: "CreatedAt"
					Labels: {}
				}
				UpdatedAt: {
					@hof(field)

					// this should be a string you can use within your templates
					Type: "datetime"
					Name: "UpdatedAt"
					Labels: {}
				}
				active: {
					@hof(field)

					// this should be a string you can use within your templates
					Type:    "bool"
					Name:    "active"
					Default: "false"
					Labels: {}
					Nullable: false
				}

				// this is the new field
				// username: fields.String
				DeletedAt: {
					@hof(field)

					// this should be a string you can use within your templates
					Type: "datetime"
					Name: "DeletedAt"
					Labels: {}
				}
			}
			Name: "User"
			Labels: {}
		}
		UserProfile: {
			@hof(model)
			Fields: {
				firstName: {
					@hof(field)

					// this should be a string you can use within your templates
					Type:     "string"
					Name:     "firstName"
					Length:   64
					Unique:   false
					Nullable: false
					Labels: {}
					Validation: {
						Max: 64
					}
				}
				middleName: {
					@hof(field)

					// this should be a string you can use within your templates
					Type:     "string"
					Name:     "middleName"
					Length:   64
					Unique:   false
					Nullable: false
					Labels: {}
					Validation: {
						Max: 64
					}
				}
				ID: {
					@hof(field)

					// this should be a string you can use within your templates
					Type:     "uuid"
					Name:     "ID"
					Nullable: false
					Unique:   true
					Default:  "gen_random_uuid()"
					Labels: {}
					Validation: {
						Format: "uuid"
					}
				}
				CID: {
					@hof(field)

					// this should be a string you can use within your templates
					Type:     "cuid"
					Name:     "CID"
					Nullable: false
					Labels: {}
					Unique: true
				}
				CreatedAt: {
					@hof(field)

					// this should be a string you can use within your templates
					Type: "datetime"
					Name: "CreatedAt"
					Labels: {}
				}
				UpdatedAt: {
					@hof(field)

					// this should be a string you can use within your templates
					Type: "datetime"
					Name: "UpdatedAt"
					Labels: {}
				}
				lastName: {
					@hof(field)

					// this should be a string you can use within your templates
					Type:     "string"
					Name:     "lastName"
					Length:   64
					Unique:   false
					Nullable: false
					Labels: {}
					Validation: {
						Max: 64
					}
				}
				DeletedAt: {
					@hof(field)

					// this should be a string you can use within your templates
					Type: "datetime"
					Name: "DeletedAt"
					Labels: {}
				}
			}
			Name: "UserProfile"
			Labels: {}
		}
	}
	Labels: {}
}