package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#DocCommand: schema.#Command & {
	TBD:   "Ø"
	Name:  "doc"
	Usage: "doc"
	Aliases: ["docs"]
	Short: "Generate and view documentation."
	Long:  Short
}

#TourCommand: schema.#Command & {
	TBD:   "Ø"
	Name:  "tour"
	Usage: "tour"
	Short: "take a tour of the hof tool"
	Long:  Short
}

#TutorialCommand: schema.#Command & {
	TBD:   "Ø"
	Name:  "tutorial"
	Usage: "tutorial"
	Short: "tutorials to help you learn hof right in hof"
	Long:  Short
}
