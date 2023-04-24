package Datamodel

msg_20230423230059: "first checkpoint"

// The models we want in our database as tables & columns
ver_20230423230059: {
	$hof: {
		datamodel: {
			root:    true
			history: true
		}
		metadata: {
			name: "Datamodel"
			id:   "datamodel"
		}
	}
	History: []
	// these are the models for the application
	// they can map onto database tables and apis
	Models: {
		// Actual Models
		User: {
			// schema for $hof: ...
			$hof: {
				apiVersion: "v1beta1"
				// typical metadata
				metadata: name: "User"

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
				ID: {
					Name:     "ID"
					Type:     "uuid"
					Nullable: false
					Unique:   true
					Default:  "gen_random_uuid()"
					Validation: Format: "uuid"
				}
				CreatedAt: {
					Name: "CreatedAt"
					Type: "datetime"
				}
				UpdatedAt: {
					Name: "UpdatedAt"
					Type: "datetime"
				}
				DeletedAt: {
					Name: "DeletedAt"
					Type: "datetime"
				}
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

				// relation fields
				Profile: {
					Name:     "Profile"
					Type:     "uuid"
					Nullable: false
					Unique:   true
					Default:  "gen_random_uuid()"
					Validation: Format: "uuid"

					// relation type, open to be flexible
					Relation: {
						Name:  "Profile"
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
		$hof: datamodel: {
			node:    true
			ordered: true
		}
		UserProfile: {
			// schema for $hof: ...
			$hof: {
				apiVersion: "v1beta1"
				// typical metadata
				metadata: name: "UserProfile"

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
	}

	// OrderedModels: [...Model] will be
	// inject here for order stability

	// these are the models for the application
	// they map onto database tables
}
