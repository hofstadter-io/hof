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
		Name:  "view"
		Usage: "view"
		Short: "view data model information"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "list"
		Usage: "list"
		Short: "list the known modelsets"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "status"
		Usage: "status"
		Short: "show the current status for a modelset"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "graph"
		Usage: "graph"
		Short: "show the relationship graph for a modelset"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "diff"
		Usage: "diff"
		Short: "show the current diff for a modelset"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "migrate"
		Usage: "migrate"
		Short: "create the next migration for a modelset"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "test"
		Usage: "test"
		Short: "test the current migration and diff for a modelset"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "delete"
		Usage: "delete"
		Short: "delete a modelset permentantly"
		Long:  Short
	}]
}
