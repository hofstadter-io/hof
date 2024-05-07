package diff

Datamodel: Models: {
	User: Fields: "+": {
		// relation fields
		Profile: {
			Name:     "Profile"
			Plural:   "Profiles"
			Type:     "uuid"
			Nullable: false
			Unique:   true
			Validation: Format: "uuid"

			// relation type, open to be flexible
			Relation: {
				Name:  "Profile"
				Type:  "has-one"
				Other: "Models.UserProfile"
			}

			// we can enrich this for various types
			// in our app or other reusable datamodels
		}
	}
	"+": UserProfile: {
		// for easy access
		Name:   "UserProfile"
		Plural: "UserProfiles"

		// These are the fields of a model
		// they can map onto database columnts and form fields
		Fields: {
			About: {
				Name:   "About"
				Plural: "Abouts"
				SQL: Type: "character varying(64)"
				Type:     "string"
				Length:   64
				Unique:   false
				Nullable: false
				Validation: Max: 64
			}
			Avatar: {
				Name:   "Avatar"
				Plural: "Avatars"
				SQL: Type: "character varying(64)"
				Type:     "string"
				Length:   64
				Unique:   false
				Nullable: false
				Validation: Max: 64
			}
			Social: {
				Name:   "Social"
				Plural: "Socials"
				SQL: Type: "character varying(64)"
				Type:     "string"
				Length:   64
				Unique:   false
				Nullable: false
				Validation: Max: 64
			}
			ID: {
				Name:     "ID"
				Plural:   "IDs"
				Type:     "uuid"
				Nullable: false
				Unique:   true
				Default:  "uuid_generate_v4()"
				Validation: Format: "uuid"
			}
			CreatedAt: {
				Name:   "CreatedAt"
				Plural: "CreatedAts"
				Type:   "datetime"
			}
			Owner: {
				Name:     "Owner"
				Plural:   "Owners"
				Type:     "uuid"
				Nullable: false
				Unique:   true
				Validation: Format: "uuid"

				// relation type, open to be flexible
				Relation: {
					Name:  "Owner"
					Type:  "belongs-to"
					Other: "Models.User"
				}

				// we can enrich this for various types
				// in our app or other reusable datamodels
			}
			UpdatedAt: {
				Name:   "UpdatedAt"
				Plural: "UpdatedAts"
				Type:   "datetime"
			}
		}

		// if we want Relations as a separate value
		// we can process the fields to extract them
		Snapshot: Timestamp: ""
		History: []
	}
}

package diff

User: Fields: "+": {
	// relation fields
	Profile: {
		Name:     "Profile"
		Plural:   "Profiles"
		Type:     "uuid"
		Nullable: false
		Unique:   true
		Validation: Format: "uuid"

		// relation type, open to be flexible
		Relation: {
			Name:  "Profile"
			Type:  "has-one"
			Other: "Models.UserProfile"
		}

		// we can enrich this for various types
		// in our app or other reusable datamodels
	}
}

