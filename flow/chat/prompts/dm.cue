package prompts

Datamodel: {
	context: #"""
		\(_mainContent)
		"""#
	examples: _examples

	_startingJSON: #"""
		Datamodel: {
			@datamodel()
			Name: "<datamodel-name>"
			Models: {
				User: {
					id: "uuid"
					created_at: "datetime"
					updated_at: "datetime"
					deleted_at: "datetime"
					email: "email"
					active: "bool"
					verified: "bool"
				}
			}
		}
		"""#

	_mainContent: #"""
		TYPES: [string, int, bool, float, uuid, datetime, email, url]
		RELATIONS: [belongs-to, has-one, has-many, many-to-many]

		// SCHEMA
		```
		Datamodel: {
			@datamodel()
			Name: "<datamodel-name>"
			Models: {
				<model-name>: {
					<field-name>: "<field-type>"
					<field-name>: "<field-type>"
					$relations: {
						<relation-name>: {
							type: "<relation-type>"
							model: "<model-name>"
						}
						<relation-name>: {
							...
						}
					}
				}
				<model-name>: {
					...
				}
			}
		}
		```

		Your task is to modify the original JSON object according to the instructions.
		The JSON should conform to the SCHEMA and follow the GUIDELINES.
		Use the following guidelines when performing the task.

		GUIDELINES:
		- The Datamodel is a loose representation of a SQL database, models are tables and fields are columns.
		- The Datamodel is composed of Models, Models are composed of fields and $relations.
		- Do NOT add extra models the user does not ask for.
		- Do NOT place a Model within another Model. You may only modify them to add $relations.
		- You are allowed to make assumptions about the need for new models or fields if the instructions seem to imply their need.
		- <field-type> must come from the TYPES list,
		- The common database fields are id, created_at, updated_at, and deleted_at.
		- When adding a new model, include the common database fields, unless instructed otherwise.
		- If the instructions do not specify the field type, you should make your best guess.
		- You should try to keep the number of models concise and not introduce unnecessary duplication of information.
		- If a field can be calculated by a SQL query on another table, don't add it.
		- <relation-type> must come from the RELATIONS list.
		- Models can have relations between them. If you make a relation, there must be a model for both sides.
		- If a user has something, this implies a new Model and Relation. It is up to you to determine the correct relation type.
		- "has-many" and "many-to-many" relations should be named as the plural of the model they refer to.
		- "many-to-many" relations require an extra model to hold the linking information.
		- Remove quotes from keys, unless they contain unusual characters.

		You should only output the results as JSON with no explanation, extra labels, or other words. (important)
		"""#
}

_examples: [{
	input: "Create a data model called Interludes. Users have a Profile with an avatar, about section, and their current status."
	output: """
		Datamodel: {
		  @datamodel()
		  Name: "Interludes"
		  Models: {
		    User: {
		      id: "uuid"
		      created_at: "datetime"
		      updated_at: "datetime"
		      deleted_at: "datetime"
		      email: "email"
		      active: "bool"
		      verified: "bool"
		      $relations: {
		        profile: {
		          type: "has-one"
		          model: "Profile"
		        }
		      }
		    }
		    Profile: {
		      id: "uuid"
		      created_at: "datetime"
		      updated_at: "datetime"
		      deleted_at: "datetime"
		      avatar: "url"
		      about: "string"
		      status: "string"
		      $relations: {
		        user: {
		          type: "belongs-to"
		          model: "User"
		        }
		      }
		    }
		  }
		}
		"""
}, {

	input: "Create a data model called Interludes. Users can have many posts. They can write them as drafts and publish them later."
	output: """
		Datamodel: {
		  @datamodel()
		  Name: "Interludes"
		  Models: {
		    User: {
		      id: "uuid"
		      created_at: "datetime"
		      updated_at: "datetime"
		      deleted_at: "datetime"
		      email: "email"
		      active: "bool"
		      verified: "bool"
		      $relations: {
		        posts: {
		          type: "has-many"
		          model: "Post"
		        }
		      }
		    }
		    Post: {
		      id: "uuid"
		      created_at: "datetime"
		      updated_at: "datetime"
		      deleted_at: "datetime"
		      title: "string"
		      body: "string"
		      published: "datetime"
		      $relations: {
		        author: {
		          type: "belongs-to"
		          model: "User"
		        }
		      }
		    }
		  }
		}
		"""
}, {
	input: """
		Create a data model called Interludes with Users.
		Users have a Profile with an avatar, about section, and their current status.
		Users can have many posts. They can write them as drafts and publish them later.
		"""

	output: """
		Datamodel: {
		  @datamodel()
		  Name: "Interludes"
		  Models: {
		    User: {
		      id: "uuid"
		      created_at: "datetime"
		      updated_at: "datetime"
		      deleted_at: "datetime"
		      email: "email"
		      active: "bool"
		      verified: "bool"
		      $relations: {
		        profile: {
		          type: "has-one"
		          model: "Profile"
		        }
		        posts: {
		          type: "has-many"
		          model: "Post"
		        }
		      }
		    }
		    Profile: {
		      id: "uuid"
		      created_at: "datetime"
		      updated_at: "datetime"
		      deleted_at: "datetime"
		      avatar: "url"
		      about: "string"
		      status: "string"
		      $relations: {
		        user: {
		          type: "belongs-to"
		          model: "User"
		        }
		      }
		    }
		    Post: {
		      id: "uuid"
		      created_at: "datetime"
		      updated_at: "datetime"
		      deleted_at: "datetime"
		      title: "string"
		      body: "string"
		      published: "datetime"
		      $relations: {
		        author: {
		          type: "belongs-to"
		          model: "User"
		        }
		      }
		    }
		  }
		}
		"""
}]
