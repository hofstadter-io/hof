package datamodel

import "github.com/hofstadter-io/hof/schema"

#Bucket: schema.#Model & {
	plural: "buckets"

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
		short: schema.#UUID
		name: schema.#String
		description: schema.#String
	}

	Index: [
		{ unique: true, fields: ["project_id", "name"] }
	]

	Lookup: {
		name: Fields.name
	}

}
