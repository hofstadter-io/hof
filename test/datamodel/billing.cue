package datamodel

import "github.com/hofstadter-io/hof/schema"

#BillingAccount: schema.#Model & {
	Permissioned: {
		object: "optional"
		group: true
		user: true
	}

	Owned: {
		by: "#Account"
		type: "has-many"
	}

	Fields: {
		name: schema.#String
		type: schema.#String
		state: schema.#String
		stripeToken: schema.#String
		stripeCustomer: schema.#String
	}

	Index: [
		{ unique: true, fields: ["account_id", "name"] }
	]

	Lookup: {
		name: Fields.name
	}

	Relations: {
		project: {
			relation: "has-many"
			type: "#Project"
		}
	}

}

#StripeSubscription: schema.#Model & {
	Permissioned: {
		object: "optional"
		group: true
		user: true
	}

	Owned: {
		by: "#Account"
		type: "has-many"
	}

	Fields: {
		name: schema.#String
		type: schema.#String
		stripeID: schema.#UUID
	}

	Index: [
		{ unique: true, fields: ["account_id", "name"] }
	]

	Lookup: {
		name: Fields.name
	}

	Relations: {
		billing: {
			relation: "belongs-to-many"
			type: "#BillingAccount"
			nullable: true
		}
	}
}
