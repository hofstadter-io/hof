package datamodel

import "github.com/hofstadter-io/hof/schema"

#Secret: schema.#Model & {
	plural: "secrets"

	Permissioned: {
		group: true
		user: true
	}

	Quotas: {
		enabled: true
	}

	Owned: {
		by: "#Project"
		type: "has-many"
	}

	Fields: schema.#CommonFields & {
		name: schema.#String
		description: schema.#String & { length: 256 }
	}

	Index: [
		{ unique: true, fields: ["project_id", "name"] }
	]

	Lookup: {
		short: Fields.short
		name: Fields.name
	}

}
