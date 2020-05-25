package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#DatamodelCommand: schema.#Command & {
	TBD:   "α"
	Name:  "datamodel"
	Usage: "datamodel"
	Aliases: ["dmod"]
	Short: "create, view, diff, calculate / migrate, and manage your data models"
	Long:  Short

	OmitRun: true

	Commands: [{
		TBD:   "α"
		Name:  "get"
		Usage: "get"
		Short: "find and display data models"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "status"
		Usage: "status"
		Short: "print the data model status"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "visualize"
		Usage: "visualize"
		Short: "visualize a data model"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "diff"
		Usage: "diff"
		Short: "show the current diff for a data model"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "history"
		Usage: "history"
		Short: "show the history for a data model"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "migrate"
		Usage: "migrate"
		Short: "calculate a changeset for a data model"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "apply"
		Usage: "apply"
		Short: "apply a migraion sequence against a data store"
		Long:  Short
	}]
}
