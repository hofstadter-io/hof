package User

msg_20230424080137: "first dm"
// Actual Models
ver_20230424080137: {
	// schema for $hof: ...
	$hof: {
		apiVersion: "v1beta1"
		// typical metadata
		metadata: {
			id:   "user"
			name: "User"
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
	Name: "User"

	// These are the fields of a model
	// they can map onto database columnts and form fields
	Fields: {
		email: {
			Name:     "email"
			sqlType:  "varchar(64)"
			Type:     "string"
			Length:   64
			Unique:   true
			Nullable: false
			Validation: {
				Max:    64
				Format: "email"
			}
		}
		password: {
			Name:     "password"
			Bcrypt:   true
			sqlType:  "varchar(64)"
			Type:     "string"
			Length:   64
			Unique:   false
			Nullable: false
			Validation: Max: 64
		}
		active: {
			Name:     "active"
			Type:     "bool"
			Default:  "false"
			Nullable: false
		}

		// this is the new field
		username: {
			Name:     "username"
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

		// relation fields
		Profile: {
			Name:     "Profile"
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
