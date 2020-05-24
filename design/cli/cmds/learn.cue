package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#TourCommand: schema.#Command & {
	TBD:   " "
	Name:  "tour"
	Usage: "tour"
	Short: "Take a tour of the hof tool"
	Long:  Short
}

#TutorialCommand: schema.#Command & {
	TBD:   " "
	Name:  "tutorial"
	Usage: "tutorial"
	Short: "Tutorials to help you learn hof right in hof"
	Long:  Short
}
