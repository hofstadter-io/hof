package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#LabelCommand: schema.#Command & {
	TBD:   "α"
	Name:  "label"
	Usage: "label"
	Aliases: ["labels", "attrs"]
	Short: "manage labels for resources and more"
	Long:  Short

	OmitRun: true

	Commands: [{
		TBD:   "α"
		Name:  "info"
		Usage: "info"
		Short: "print info about labels in your workspace or system"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "add"
		Usage: "add"
		Short: "add labels to your workspace or system"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "edit"
		Usage: "edit"
		Short: "edit labels in your workspace or system configurations"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "delete"
		Usage: "delete"
		Short: "delete labels from your workspace or system"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "get"
		Usage: "get"
		Short: "find and display labels on resources"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "set"
		Usage: "set"
		Short: "find and apply labels to resources"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "remove"
		Usage: "remove"
		Short: "find and remove labels from resources"
		Long:  Short
	}]
}
