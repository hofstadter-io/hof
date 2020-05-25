package schema

import (
	"github.com/hofstadter-io/hof/schema/datamodel"
)

#Datamodel: [N=string]: #Datamodel & { Name: N, ... }
#Datamodel: {
  Name: string

	Workbase: string | *"" // usually current directory
	Entrypoint: string     // paths from workbase

	Modelsets: datamodel.#Modelsets
	Models: #Models
	Views: #Models

	...
}
