// Traditional database model which maps onto tables & columns.
Datamodel: {
	// schema for #hof: ...
	#hof: {
		// #hof version
		apiVersion: "v1beta1"

		// typical metadata
		metadata: {}

		// hof/datamodel
		datamodel: {
			// define the root of a datamodel
			root: true

			// instruct history to be tracked
			history: true

			// instruct ordrered version of the fields
			// to be injected as a peer value
			ordered: false

			// tell hof this is a node of interest for
			// the inspection commands (list,info)
			node: false

			// tell hof to track this as a raw CUE value
			// (partially implemented)
			cue: false
		}
	}
	Snapshot: {
		Timestamp: ""
	}

	// these are the models for the application
	// they can map onto database tables and apis
	Models: {
		User: {
			// for easy access
			Name:   "User"
			Plural: "Users"

			// These are the fields of a model
			// they can map onto database columnts and form fields
			Fields: {
				ID: {
					Name:     "ID"
					Plural:   "IDs"
					Type:     "uuid"
					Nullable: false
					Unique:   true
					Validation: {
						Format: "uuid"
					}
					#hof: {
						metadata: {
							name: "ID"
						}
					}
				}
				CreatedAt: {
					Name:   "CreatedAt"
					Plural: "CreatedAts"
					Type:   "datetime"
					#hof: {
						metadata: {
							name: "CreatedAt"
						}
					}
				}
				UpdatedAt: {
					Name:   "UpdatedAt"
					Plural: "UpdatedAts"
					Type:   "datetime"
					#hof: {
						metadata: {
							name: "UpdatedAt"
						}
					}
				}
				DeletedAt: {
					Name:   "DeletedAt"
					Plural: "DeletedAts"
					Type:   "datetime"
					#hof: {
						metadata: {
							name: "DeletedAt"
						}
					}
				}
				email: {
					Name:     "email"
					Plural:   "emails"
					Type:     "string"
					Length:   64
					Unique:   true
					Nullable: false
					Validation: {
						Max:    64
						Format: "email"
					}
					#hof: {
						metadata: {
							name: "email"
						}
					}
				}
				username: {
					Name:     "username"
					Plural:   "usernames"
					Type:     "string"
					Length:   64
					Unique:   false
					Nullable: false
					Validation: {
						Max: 64
					}
					#hof: {
						metadata: {
							name: "username"
						}
					}
				}
				password: {
					Name:     "password"
					Plural:   "passwords"
					Bcrypt:   true
					Type:     "string"
					Length:   64
					Unique:   false
					Nullable: false
					Validation: {
						Max: 64
					}
					#hof: {
						metadata: {
							name: "password"
						}
					}
				}
				verified: {
					Name:     "verified"
					Plural:   "verifieds"
					Type:     "bool"
					Default:  "false"
					Nullable: false
					#hof: {
						metadata: {
							name: "verified"
						}
					}
				}
				active: {
					Name:     "active"
					Plural:   "actives"
					Type:     "bool"
					Default:  "false"
					Nullable: false
					#hof: {
						metadata: {
							name: "active"
						}
					}
				}
				persona: {
					Name:   "persona"
					Plural: "personas"
					Type:   "string"
					Vals: ["guest", "user", "admin", "owner"]
					Nullable: false
					Default:  "user"
					#hof: {
						metadata: {
							name: "persona"
						}
					}
				}
				#hof: {
					datamodel: {
						node:    true
						ordered: true
					}
				}
			}

			// if we want Relations as a separate value
			// we can process the fields to extract them
			// schema for #hof: ...
			#hof: {
				// #hof version
				apiVersion: "v1beta1"

				// typical metadata
				metadata: {
					name: "User"
				}

				// hof/datamodel
				datamodel: {
					// define the root of a datamodel
					root: false

					// instruct history to be tracked
					history: true

					// instruct ordrered version of the fields
					// to be injected as a peer value
					ordered: false

					// tell hof this is a node of interest for
					// the inspection commands (list,info)
					node: false

					// tell hof to track this as a raw CUE value
					// (partially implemented)
					cue: false
				}
			}
			Snapshot: {
				Timestamp: ""
			}
			History: []
		}
		#hof: {
			datamodel: {
				node:    true
				ordered: true
			}
		}
	}

	// OrderedModels: [...Model] will be
	// inject here for order stability
	History: []
}
