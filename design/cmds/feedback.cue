package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#FeedbackCommand: schema.#Command & {
	// TBD:   "Ø"
	Name:  "feedback"
	Usage: "feedback <message>"
	Aliases: ["hi"]
	Short: "send feedback, bug reports, or any message"
	Long: """
		Opens an issue on GitHub with some fields prefilled out
		"""
}
