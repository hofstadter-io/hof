package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

FeedbackCommand: schema.Command & {
	// TBD:   "Ø"
	Name:  "feedback"
	Usage: "feedback <message>"
	Aliases: ["hi", "ask", "report"]
	Short: "open an issue or discussion on GitHub"
	Long: """
		Opens an issue or discussion on GitHub with some fields prefilled out
		"""

	Pflags: [{
		Name:    "issue"
		Long:    "issue"
		Short:   "I"
		Type:    "bool"
		Default: "false"
		Help:    "create an issue (discussion is default)"
	}, {
		Name:    "labels"
		Long:    "labels"
		Short:   "L"
		Type:    "string"
		Default: "\"feedback\""
		Help:    "labels,comma,separated"
	}]
}
