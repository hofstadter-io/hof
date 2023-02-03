package main

import (
	"github.com/hofstadter-io/hof/schema/dm"
)

// example using hof's common datamodel
Datamodel: dm.Datamodel & {
	Models: {
		User: {
			Fields: {
				name: Type: "string"
				email: Type: "string"
				active: Type: "bool"
			}
			Relations: {
				Profile: {
					Type: "has-one"
					Path: "Models.Profile"
				}
			}
		}

		Profile: {
			Fields: {
				about: Type: "string"
			}
		}
	}
}

// This is to try out a CUE object in datamodel
Config: dm.Value & {
	#hof: metadata: name: "Config"

	Server: {
		@node()
		Host: string
		Port: string
	}

	Runtime: {
		@node()
		Debug: bool
		LogLevel: string
	}

	Auth: {
		@node()
		Password: bool
		Apikey:   bool
		OAuth: {
			github: bool
			google: bool
			apple:  bool
		}
	}

	Storage: {
		@node()

		Database: {
			@node()
			name:   string
			engine: string
			config: _
		}

		Bucket: {
			@node()
			name:   string
			engine: string
			config: _
		}

		Cache: {
			@node()
			name:   string
			engine: string
			config: _
		}
	}
}
