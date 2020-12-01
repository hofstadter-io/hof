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
	Long:  Short

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
		Name:  "create"
		Usage: "create"
		Aliases: ["c"]
		Short: "create a new datamodel"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "get"
		Usage: "get"
		Aliases: ["g"]
		Short: "find and display data models"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "edit"
		Usage: "edit"
		Aliases: ["e"]
		Short: "find and edit data models"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "delete"
		Usage: "delete"
		Aliases: ["del"]
		Short: "find and delete data models"
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
		Name:  "visualize"
		Usage: "visualize"
		Aliases: ["v", "viz", "show"]
		Short: "visualize a data model"
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
		Name:  "migrate"
		Usage: "migrate"
		Aliases: ["mig", "migs", "migrations"]
		Short: "calculate a changeset for a data model"
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
