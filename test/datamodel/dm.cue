package datamodel

import "github.com/hofstadter-io/hof/schema/dm"

#MyModels: dm.#Datamodel & {
	Name: "MyModels"

	Models: {
		#UserModels
	}
}
