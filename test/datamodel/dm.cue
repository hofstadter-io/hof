package datamodel

import "github.com/hofstadter-io/hof/schema/dm"

#MyModels: dm.#Datamodel & {
	Name: "MyModels"

	Models: [
		for M in #UserModels { M },
	]
}
