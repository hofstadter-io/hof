package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#ModelsetCommand: schema.#Command & {
	TBD:  "α"
	Name:  "modelset"
	Usage: "modelset"
	Aliases: ["mset"]
	Short: "create, view, migrate, and understand your modelsets."
	Long:  Short

	OmitRun: true

	Commands: [{
		TBD:  "α"
		Name:  "create"
		Usage: "create"
		Short: "create a modelset"
		Long:  Short

		Args: [{
			Name:     "name"
			Type:     "string"
			Required: true
			Help:     "modelset name"
		}, {
			Name:    "entrypoint"
			Type:    "string"
			Default: "\"models\""
			Help:    "the directory where your modelset will exist"
		}]
	}, {
		TBD:  "α"
		Name:  "view"
		Usage: "view"
		Short: "view modelset information"
		Long:  Short
		Args: [{
			Name:     "name"
			Type:     "string"
			Required: true
			Help:     "modelset name"
		}]
	}, {
		TBD:  "α"
		Name:  "list"
		Usage: "list"
		Short: "list the known modelsets"
		Long:  Short
	}, {
		TBD:  "α"
		Name:  "status"
		Usage: "status"
		Short: "show the current status for a modelset"
		Long:  Short
		Args: [{
			Name:     "name"
			Type:     "string"
			Required: true
			Help:     "modelset name"
		}]
	}, {
		TBD:  "α"
		Name:  "graph"
		Usage: "graph"
		Short: "show the relationship graph for a modelset"
		Long:  Short
		Args: [{
			Name:     "name"
			Type:     "string"
			Required: true
			Help:     "modelset name"
		}]
	}, {
		TBD:  "α"
		Name:  "diff"
		Usage: "diff"
		Short: "show the current diff for a modelset"
		Long:  Short
		Args: [{
			Name:     "name"
			Type:     "string"
			Required: true
			Help:     "modelset name"
		}]
	}, {
		TBD:  "α"
		Name:  "migrate"
		Usage: "migrate"
		Short: "create the next migration for a modelset"
		Long:  Short
		Args: [{
			Name:     "name"
			Type:     "string"
			Required: true
			Help:     "modelset name"
		}]
	}, {
		TBD:  "α"
		Name:  "test"
		Usage: "test"
		Short: "test the current migration and diff for a modelset"
		Long:  Short
		Args: [{
			Name:     "name"
			Type:     "string"
			Required: true
			Help:     "modelset name"
		}]
	}, {
		TBD:  "α"
		Name:  "delete"
		Usage: "delete"
		Short: "delete a modelset permentantly"
		Long:  Short
		Args: [{
			Name:     "name"
			Type:     "string"
			Required: true
			Help:     "modelset name"
		}]
	}]
}
