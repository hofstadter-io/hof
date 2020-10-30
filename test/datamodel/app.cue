package datamodel

import "github.com/hofstadter-io/hof/schema"

#App: schema.#Model & {
	plural: "apps"

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
		short: schema.#XUID
		name: schema.#String
		type: schema.#String
		description: schema.#String
		version: schema.#String
	}

	Relations: {
		appDeployment: {
			relation: "has-many"
			type: "#AppDeployment"
		}
	}
	
	Index: [
		{ unique: true, fields: ["project_id", "name"] }
	]

	Lookup: {
		short: Fields.short
		name: Fields.name
	}

}

#AppDeployment: schema.#Model & {
	plural: "app-deployments"

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
		short: schema.#XUID
		name: schema.#String
		mode: schema.#String
		version: schema.#String
	}

	Relations: {
		app: {
			relation: "belongs-to-one"
			type: "#App"
		}
	}

	Index: [
		{ unique: true, fields: ["project_id", "name"] }
	]

	Lookup: {
		short: Fields.short
		name: Fields.name
	}

}
