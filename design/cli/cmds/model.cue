package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#ModelsetCommand: schema.#Command & {
	Name:  "modelset"
	Usage: "modelset"
	Aliases: ["mset"]
	Short: "create, view, migrate, and understand your modelsets."
	Long:  Short

	OmitRun: true

	Commands: [{
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
		Name:  "list"
		Usage: "list"
		Short: "list the known modelsets"
		Long:  Short
	}, {
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
