package datamodel

import "github.com/hofstadter-io/hof/schema"

#Account: schema.#Model & {

	Permissioned: {
		Object: "required"
		Group: true
		User: true
	}

	Fields: schema.#CommonFields & {
		short: schema.#XUID
		name: schema.#String
		state: schema.#String & { length: 32 }
	}

	Lookup: {
		short: Fields.short
		name: Fields.name
	}
}
