package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#ChatCommand: schema.#Command & {
	Name:  "chat"
	Usage: "chat [args]"
	Short: "Co-design with AI (alpha)"
	Long:  #ChatRootHelp

	Flags: [...schema.#Flag] & [ {
		Name:    "model"
		Type:    "string"
		Default: "\"gpt-3.5-turbo\""
		Help:    "LLM model to use [gpt-3.5-turbo,gpt-4]"
		Long:    "model"
		Short:   "M"
	},
		{
			Name:    "prompt"
			Type:    "string"
			Default: "\"\""
			Help:    "path to the system prompt, the first message in the chat"
			Long:    "prompt"
			Short:   "P"
		},
		{
			Name:    "outfile"
			Type:    "string"
			Default: "\"\""
			Help:    "path to write the output to"
			Long:    "outfile"
			Short:   "O"
		},
	]
}

#ChatRootHelp: #"""
	Use chat to work with hof features or from modules you import.
	Module authors can provide custom prompts for their schemas.

	Currently, only ChatGPT is supported. You can use any of the
	gpt-3.5 or gpt-4 models. The flag should match OpenAI API options.

	Set OPENAI_API_KEY
	"""#
