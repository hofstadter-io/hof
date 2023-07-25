package schema

import (
	"github.com/hofstadter-io/hof/schema/dm"
	"github.com/hofstadter-io/hof/schema/gen"
)

// The schema for your thing
MyThing: {
	// with any fields you want
	// these are typically specific to your application
}

// A hof datamodel, also part of your schema
Datamodel: dm.Datamodel

// The generator schema your users will use
//   it is based on the core hof/gen.Generator
Generator: gen.Generator & {

	// User-facing inputs to your generator
	"MyThing":   MyThing
	"Datamodel": Datamodel

	...
}
