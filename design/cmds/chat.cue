package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#ChatCommand: schema.#Command & {
	Name:  "chat"
	Usage: "chat jsonfile [extra args]"
	Short: "Co-design with AI (alpha)"
	Long:  #ChatRootHelp

	Args: [{
		Name:     "jsonfile"
		Type:     "string"
		Required: true
		Help:     "The JSON file to operate on"
	}, {
		Name:     "instructions"
		Type:     "string"
		Required: true
		Help:     "The instructions for the AI"
	}, {
		Name: "extra"
		Type: "[]string"
		Rest: true
		Help: "extra arguments for the chat, tbd what they are"
	}]

	Flags: [...schema.#Flag] & [ {
		Name:    "model"
		Type:    "string"
		Default: "\"gpt-3.5-turbo\""
		Help:    "LLM model to use [gpt-3.5-turbo,gpt-4]"
		Long:    "model"
		Short:   "M"
	},
		{
			Name:    "parameters"
			Type:    "string"
			Default: "\"\""
			Help:    "path to a config file containing model parameters"
			Long:    "parameters"
			Short:   "P"
		},
		{
			Name:    "generator"
			Type:    "[]string"
			Default: "nil"
			Help:    "generator tags to run, default is all"
			Long:    "generator"
			Short:   "G"
		}]
}

#ChatRootHelp: #"""
	hof chat synergizes LLM technology with our code gen technology.
	You can chat with the AI to help you write designs and code.
	Module authors can provide custom prompts for their schemas.

	Currently, only ChatGPT is supported. You can use any of the
	gpt-3.5 or gpt-4 models. The flag should match OpenAI API options.

	Set OPENAI_API_KEY
	"""#
