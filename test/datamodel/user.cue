package datamodel

import "github.com/hofstadter-io/hof/schema/dm"

#UserModels: dm.#Models & {
	"User": #User
	"UserProfile": #UserProfile
}

#User: dm.#Model & {
	Fields: {
		dm.#CommonFields
		email: dm.#Email	
		persona: dm.#Enum & {
			vals: ["guest", "user", "admin", "owner"]
			default: "guest"
		}	
		password: dm.#Password
		active: dm.#Bool
	}
}

#UserProfile: dm.#Model & {
	Fields: {
		dm.#CommonFields
		firstName: dm.#String
		middleName: dm.#String
		lastName: dm.#String
	}
}
