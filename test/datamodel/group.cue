package datamodel

import "github.com/hofstadter-io/hof/schema"

#Group: schema.#Model & {
	Fields: {
		name: schema.#String & { nullable: false }
		short: schema.#XUID
		persona: schema.#String & { length: 32 }
		state: schema.#String & { length: 32 }
	}
}
