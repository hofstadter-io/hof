package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#GenCommand: schema.#Command & {
	// TBD:   "✓"
	Name:  "gen"
	Usage: "gen [files...]"
	Aliases: ["G"]
	Short: "generate code, data, and config from your data models and designs"
	Long: """
    generate all the things, from code to data to config...
  """

	Flags: [...schema.#Flag] & [
		{
			Name:    "stats"
			Type:    "bool"
			Default: "false"
			Help:    "Print generator statistics"
			Long:    "stats"
			Short:   "s"
		},
		{
			Name:    "generator"
			Type:    "[]string"
			Default: "nil"
			Help:    "Generators to run, default is all discovered"
			Long:    "generator"
			Short:   "g"
		},
		//{
			//Name:    "f"
			//Type:    "[]string"
			//Default: "nil"
			//Help:    "File globs to render, default is all discovered"
			//Long:    "files"
			//Short:   "f"
		//},
	]

}

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
	Name:   "geb"
	Usage:  "_geb"
	Short:  ""
	Long:   ""
	Hidden: true
	OmitTests: true
}

#LogoCommand: schema.#Command & {
	Name:   "logo"
	Usage:  "_"
	Short:  ""
	Long:   ""
	Hidden: true
	OmitTests: true
}
