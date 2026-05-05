package main

import "github.com/hofstadter-io/hof/schema/dm"

MyModels: dm.Datamodel & {
	foo: "string"
	...
}
