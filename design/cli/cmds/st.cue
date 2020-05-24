package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#StLong: """
Structural diff, merge, mask, pick, and query helpers for Cue

Commands generally have the form: <method> <op1> <op2> [...entrypoints]

Where <op> can be:
- Cue: expr: as: string
- @filename.cue: Cue: expr: as: string

If entrypoints are supplied, then an <op> without an @filename.cue will lookup from the entrypoints.
Otherwise, the <op> is interpreted as a complete Cue value.
"""

#StCommand: schema.#Command & {
	Name:  "st"
	Usage: "st"
	Aliases: ["structural"]
	Short: "Structural diff, merge, mask, pick, and query helpers for Cue"
	Long:  #StLong

	OmitRun: true

	Commands: [{
		Name:  "diff"
		Usage: "diff <orig> <next> [...entrypoints]"
		Short: "Calculate the difference between two Cue values"
		Long: """
		Calculate the difference between two Cue values
		"""

		Args: [{
			Name:     "orig"
			Type:     "string"
			Required: true
			Help:     "original value to calc diff from, see 'hof st --help' for format"
		}, {
			Name:     "next"
			Type:     "string"
			Required: true
			Help:     "next value to calc diff to, see 'hof st --help' for format"
		}, {
			Name: "entrypoints"
			Type: "[]string"
			Rest: true
			Help: "Cue entrypoints"
		}]
	}, {
		Name:  "merge"
		Usage: "merge <orig> <update> [...entrypoints]"
		Short: "Merge <new> onto <orig>, replacing values and adding new ones"
		Long:  Short

		Args: [{
			Name:     "orig"
			Type:     "string"
			Required: true
			Help:     "original value to update, see 'hof st --help' for format"
		}, {
			Name:     "update"
			Type:     "string"
			Required: true
			Help:     "update value to layer on, see 'hof st --help' for format"
		}, {
			Name: "entrypoints"
			Type: "[]string"
			Rest: true
			Help: "Cue entrypoints"
		}]
	}, {
		Name:  "pick"
		Usage: "pick <orig> <what> [...entrypoints]"
		Short: "Pick <what> Cue value(s) from <orig>"
		Long:  Short

		Args: [{
			Name:     "orig"
			Type:     "string"
			Required: true
			Help:     "original value to pick value(s) from, see 'hof st --help' for format"
		}, {
			Name:     "pick"
			Type:     "string"
			Required: true
			Help:     "the Cue value shape to extract from orig, see 'hof st --help' for format"
		}, {
			Name: "entrypoints"
			Type: "[]string"
			Rest: true
			Help: "Cue entrypoints"
		}]
	}, {
		Name:  "mask"
		Usage: "mask <orig> <what> [...entrypoints]"
		Short: "Mask <what> Cue value(s) from <orig>, thereby 'filtering' the original"
		Long:  Short

		Args: [{
			Name:     "orig"
			Type:     "string"
			Required: true
			Help:     "original value to filter with mask, see 'hof st --help' for format"
		}, {
			Name:     "what"
			Type:     "string"
			Required: true
			Help:     "the mask to filter orig with, see 'hof st --help' for format"
		}, {
			Name: "entrypoints"
			Type: "[]string"
			Rest: true
			Help: "Cue entrypoints"
		}]
	}, {
		Name:  "query"
		Usage: "query <orig> <expr> [...entrypoints]"
		Short: "Query for values matching an expr and/or attributes"
		Long:  Short

		Args: [{
			Name:     "orig"
			Type:     "string"
			Required: true
			Help:     "Cue value to query values from, see 'hof st --help' for format"
		}, {
			Name:     "expr"
			Type:     "string"
			Required: true
			Help:     "Cue expression to query with, see 'hof st --help' for format"
		}, {
			Name: "entrypoints"
			Type: "[]string"
			Rest: true
			Help: "Cue entrypoints"
		}]
	}]
}
