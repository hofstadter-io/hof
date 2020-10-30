package datamodel

import "github.com/hofstadter-io/hof/schema"

#Apikey: schema.#Model & {
	Fields: schema.#CommonFields & {
		apikey: schema.#UUID & { bcrypt: true }
		name: schema.#String
	}
}
