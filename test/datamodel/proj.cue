package datamodel

import "github.com/hofstadter-io/hof/schema"

#Project: schema.#Model & {
	Permissioned: {
		object: "optional"
		group: true
		user: true
	}

	Fields: schema.#CommonFields & {
		short: schema.#XUID
		name: schema.#String
		type: schema.#String
		billed: schema.#Bool
		state: schema.#String
	}

	Index: [
		{ unique: true, fields: ["account_id", "name"] }
	]

	Lookup: {
		name: Fields.name
		short: Fields.short
	}

	Relations: {
		billing: {
			relation: "belongs-to-one"
			type: "#BillingAccount"
			nullable: true
		}
	}
}
