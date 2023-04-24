package UserProfile

msg_20230424080345: "add posts"
ver_20230424080345: {
	// schema for $hof: ...
	$hof: {
		apiVersion: "v1beta1"
		// typical metadata
		metadata: {
			id:   "user-profile"
			name: "UserProfile"
		}

		// hof/datamodel
		datamodel: {
			root:    false
			history: true
			ordered: false
			node:    false
			cue:     false
		}
	}
	History: []

	// for easy access
	Name: "UserProfile"

	// These are the fields of a model
	// they can map onto database columnts and form fields
	Fields: {
		About: {
			Name:     "About"
			sqlType:  "varchar(64)"
			Type:     "string"
			Length:   64
			Unique:   false
			Nullable: false
			Validation: Max: 64
		}
		Avatar: {
			Name:     "Avatar"
			sqlType:  "varchar(64)"
			Type:     "string"
			Length:   64
			Unique:   false
			Nullable: false
			Validation: Max: 64
		}
		Social: {
			Name:     "Social"
			sqlType:  "varchar(64)"
			Type:     "string"
			Length:   64
			Unique:   false
			Nullable: false
			Validation: Max: 64
		}
		ID: {
			Name:     "ID"
			Type:     "uuid"
			Nullable: false
			Unique:   true
			Default:  "uuid_generate_v4()"
			Validation: Format: "uuid"
		}
		CreatedAt: {
			Name: "CreatedAt"
			Type: "datetime"
		}
		Owner: {
			Name:     "Owner"
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
			Name: "UpdatedAt"
			Type: "datetime"
		}
		$hof: datamodel: {
			node:    true
			ordered: true
		}
	}

	// if we want Relations as a separate value
	// we can process the fields to extract them
}
