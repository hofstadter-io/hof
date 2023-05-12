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

	This is an alpha stage command, expect big changes next release.
	We currently use t

	Currently, only ChatGPT is supported. You can use any of the
	gpt-3.5 or gpt-4 models. The flag should match OpenAI API options.
	While we are using the chat models, we do not support interactive yet.

	Set OPENAI_API_KEY as an environment variable.

	Examples:

	#
	# Talk to ChatGPT
	#

	# Ask of ChatGPT from strings, files, and/or stdin
	hof chat "Ask ChatGPT any question"    # as a string
	hof chat question.txt                  # from a file
	cat question.txt | hof chat -          # from stdin
	hof chat context.txt "and a question"  # mix all three

	# Provide a system message, these are special to ChatGPT
	hof chat -P prompt.txt "now answer me this..."

	# Get file embeddings
	hof chat embed file1.txt file2.txt -O embeddings.json

	#
	# Talk to your data model, this uses a special system message
	#

	# hof will use dm.cue by default
	hof chat dm "Create a data model called Interludes"
	hof chat dm "Users should have a Profile with status and about fields."

	# pass in a file to talk to a specific data model
	hof chat dm my-dm.cue "Add a Post model and make it so Users have many."
	"""#
