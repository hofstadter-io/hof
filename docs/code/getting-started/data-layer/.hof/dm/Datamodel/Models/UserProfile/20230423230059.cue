package UserProfile

msg_20230423230059: "first checkpoint"
ver_20230423230059: {
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
		Owner: {
			Name:     "Owner"
			Type:     "uuid"
			Nullable: false
			Unique:   true
			Default:  "gen_random_uuid()"
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
		$hof: datamodel: {
			node:    true
			ordered: true
		}
	}

	// if we want Relations as a separate value
	// we can process the fields to extract them
}
