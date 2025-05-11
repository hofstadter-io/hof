package datamodel

import (
	"github.com/hofstadter-io/schemas/dm"
)

MyObject: dm.Object & {

	foo: "bar"
	ans: 42

	animals: {
		cow: "moo"
		cat: "meow"
		dog: "woof"
	}
}
