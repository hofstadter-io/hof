package datamodel

import "github.com/hofstadter-io/hof/schema"

#BaseModelset: schema.#Modelset & {
	Name: "BaseModelset"

	Models: {
		User: #User
		UserProfile: #UserProfile
	}
}
