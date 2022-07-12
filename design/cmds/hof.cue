package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#FeedbackCommand: schema.#Command & {
	// TBD:   "Ø"
	Name:  "feedback"
	Usage: "feedback [email] <message>"
	Aliases: ["hi", "say", "from", "bug", "yo", "hello", "greetings", "support"]
	Short: "send feedback, bug reports, or any message"
	Long: """
		send feedback, bug reports, or any message
			email:     (optional) your email, if you'd like us to reply
			message:   your message, please be respectful to the person receiving it
		"""
}

#GebCommand: schema.#Command & {
	Name:      "geb"
	Usage:     "_geb"
	Short:     ""
	Long:      ""
	Hidden:    true
	OmitTests: true
}

#LogoCommand: schema.#Command & {
	Name:      "logo"
	Usage:     "_"
	Short:     ""
	Long:      ""
	Hidden:    true
	OmitTests: true
}
