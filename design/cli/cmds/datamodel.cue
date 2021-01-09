package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#DatamodelCommand: schema.#Command & {
	TBD:   "α"
	Name:  "datamodel"
	Usage: "datamodel"
	Aliases: ["dmod", "dm"]
	Short: "create, view, diff, calculate / migrate, and manage your data models"
	Long:  #DatamodelRootHelp

	OmitRun: true

	Pflags: [...schema.#Flag] & [
		{
			Name:    "Datamodels"
			Long:    "datamodel"
			Short:   "D"
			Type:    "[]string"
			Default: "nil"
			Help:    "Datamodels for the datamodel commands"
		},
		{
			Name:    "modelsets"
			Long:    "modelset"
			Short:   "M"
			Type:    "[]string"
			Default: "nil"
			Help:    "Modelsets for the datamodel commands"
		},
		{
			Name:    "models"
			Long:    "model"
			Short:   "m"
			Type:    "[]string"
			Default: "nil"

			Help:    "Models for the datamodel commands"
		},
	]

	Commands: [{
		TBD:   "α"
		Name:  "get"
		Usage: "get"
		Aliases: ["g"]
		Short: "find and display data models"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "status"
		Usage: "status"
		Aliases: ["st"]
		Short: "print the data model status"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "diff"
		Usage: "diff"
		Aliases: ["d"]
		Short: "show the current diff for a data model"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "history"
		Usage: "history"
		Aliases: ["hist", "h", "log", "l"]
		Short: "show the history for a data model"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "checkpoint"
		Usage: "checkpoint"
		Aliases: ["cp"]
		Short: "calculate a migration changeset for a data model"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "apply"
		Usage: "apply"
		Aliases: ["a"]
		Short: "apply a migraion sequence against a data store"
		Long:  Short
	}]
}

#DatamodelRootHelp: """
Data models are sets of models which are used in many hof processes and modules.

At their core, they represent the most abstract representation for objects and
their relations in your applications. They are extended and annotated to add
context fot their usage in different code generators: (DB vs Server vs Client).

Beyond representing models in their current form, a history is maintained so that:
  - database migrations can be created and managed
  - servers can handle multiple model versions
  - clients can implement feature flags
Much of this is actually handled by code generators and must be implemented there.
Hof handles the core data model definitions, history, and snapshot creation.
"""
