package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// Name arg
#nameArg: schema.#Arg & {
	Name:     "name"
	Type:     "string"
	Required: true
	Help:     "A name from /[a-zA-Z][a-zA-Z0-9_]*"
}

// Identifyier (name or id)
#identArg: schema.#Arg & {
	Name:     "ident"
	Type:     "string"
	Required: true
	Help:     "A name or id"
}

// input arg
#inputArg: schema.#Arg & {
	Name:     "input"
	Type:     "string"
	Required: true
}

#contextArg: schema.#Arg & {
	Name: "context"
	Type: "string"
	Help: "The hof auth context name"
}

// email for user / service account
#identityArg: schema.#Arg & {
	Name: "identity"
	Type: "string"
	Help: "A Hofstadter Studios user or service account"
}

// Studios account
#accountArg: schema.#Arg & {
	Name: "account"
	Type: "string"
	Help: "The name or id of a Hofstadter Studios account"
}

// Studios Project
#projectArg: schema.#Arg & {
	Name: "project"
	Type: "string"
	Help: "The name or id of a Hofstadter Studios project"
}

// Studios API Key
#apikeyArg: schema.#Arg & {
	Name: "apikey"
	Type: "string"
	Help: "Hofstadter Studios API Key"
}
