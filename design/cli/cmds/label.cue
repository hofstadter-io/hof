package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#LabelCommand: schema.#Command & {
	TBD:   "α"
	Name:  "label"
	Usage: "label"
	Aliases: ["l", "labels", "attrs"]
	Short: "manage labels for resources and more"
	Long:  Short

	OmitRun: true

	Commands: [{
		TBD:   "α"
		Name:  "info"
		Usage: "info"
		Aliases: ["i"]
		Short: "print info about labels in your workspace or system"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "create"
		Usage: "create"
		Aliases: ["c"]
		Short: "add labels to your workspace or system"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "get"
		Usage: "get"
		Aliases: ["g"]
		Short: "find and display labels from your workspace"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "set"
		Usage: "set"
		Aliases: ["s"]
		Short: "find and configure labels from your workspace"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "edit"
		Usage: "edit"
		Aliases: ["e"]
		Short: "edit labels in your workspace or system configurations"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "delete"
		Usage: "delete"
		Aliases: ["del", "remove"]
		Short: "delete labels from your workspace or system"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "apply"
		Usage: "apply"
		Aliases: ["a"]
		Short: "find and apply labels to resources"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "remove"
		Usage: "remove"
		Aliases: ["r"]
		Short: "find and remove labels from resources"
		Long:  Short
	}]
}

#LabelsetCommand: schema.#Command & {
	TBD:   "α"
	Name:  "labelset"
	Usage: "labelset"
	Aliases: ["L", "lset"]
	Short: "group resources, datamodels, labelsets, and more"
	Long:  Short

	OmitRun: true

	Commands: [{
		TBD:   "α"
		Name:  "info"
		Usage: "info"
		Aliases: ["i"]
		Short: "print info about labelsets in your workspace or system"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "create"
		Usage: "create"
		Aliases: ["c"]
		Short: "add labelsets to your workspace or system"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "get"
		Usage: "get"
		Aliases: ["g"]
		Short: "find and display labelsets from your workspace"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "set"
		Usage: "set"
		Aliases: ["s"]
		Short: "find and configure labelsets from your workspace"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "edit"
		Usage: "edit"
		Aliases: ["e"]
		Short: "edit labelsets in your workspace or system configurations"
		Long:  Short
	}, {
		TBD:   "α"
		Name:  "delete"
		Usage: "delete"
		Aliases: ["del"]
		Short: "delete labelsets from your workspace or system"
		Long:  Short
	}]
}
