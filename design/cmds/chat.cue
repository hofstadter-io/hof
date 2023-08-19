package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

ChatCommand: schema.Command & {
	Name:  "chat"
	Usage: "chat [args]"
	Short: "co-create with AI (alpha)"
	Long:  ChatRootHelp

	Flags: [...schema.Flag] & [{
		Name:    "model"
		Type:    "string"
		Default: "\"gpt-3.5-turbo\""
		Help:    "LLM model to use [gpt-3.5-turbo,gpt-4,bard,chat-bison]"
		Long:    "model"
		Short:   "M"
	}, {
		Name:    "system"
		Type:    "[]string"
		Default: "nil"
		Help:    "string or path to the system prompt for the LLM, concatenated"
		Long:    "system"
		Short:   "s"
	}, {
		Name:    "messages"
		Type:    "[]string"
		Default: "nil"
		Help:    "string or path to a message for the LLM"
		Long:    "message"
		Short:   "m"
	}, {
		Name:    "examples"
		Type:    "[]string"
		Default: "nil"
		Help:    "string or path to an example pair for the LLM"
		Long:    "example"
		Short:   "e"
	}, {
		Name:    "outfile"
		Type:    "string"
		Default: "\"\""
		Help:    "path to write the output to"
		Long:    "outfile"
		Short:   "O"
	}, {
		Name:    "Choices"
		Type:    "int"
		Default: "1"
		Help:    "param: choices or N (openai)"
		Long:    "choices"
		Short:   "N"
	}, {
		Name:    "MaxTokens"
		Type:    "int"
		Default: "256"
		Help:    "param: MaxTokens"
		Long:    "max-tokens"
	}, {
		Name:    "temperature"
		Type:    "float64"
		Default: "0.8"
		Help:    "param: temperature"
		Long:    "temp"
	}, {
		Name:    "TopP"
		Type:    "float64"
		Default: "0.42"
		Help:    "param: TopP"
		Long:    "topp"
	}, {
		Name:    "TopK"
		Type:    "int"
		Default: "40"
		Help:    "param: TopK (google)"
		Long:    "topk"
	}, {
		Name:    "Stop"
		Type:    "[]string"
		Default: "nil"
		Help:    "param: Stop (openai)"
		Long:    "stop"
	}]

	Commands: [{
		Name:  "list"
		Usage: "list"
		Short: "print available chat plugins in the current module"
		Long:  Short
	}, {
		Name:  "with"
		Usage: "with [name]"
		Short: "chat with a plugin in the current module"
		Long:  Short
		Args: [{
			Name:     "name"
			Type:     "string"
			Required: true
			Help:     "name of the chat plugin to display details for"
		}, {
			Name: "entrypoints"
			Type: "[]string"
			Rest: true
			Help: "CUE entrypoints like most other commands"
		}]
	}, {
		Name:  "info"
		Usage: "info [name]"
		Short: "print details of a specific chat plugin"
		Long:  Short
		Args: [{
			Name:     "name"
			Type:     "string"
			Required: true
			Help:     "name of the chat plugin to display details for"
		}, {
			Name: "entrypoints"
			Type: "[]string"
			Rest: true
			Help: "CUE entrypoints like most other commands"
		}]
	}]
}

ChatRootHelp: #"""
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
	# Talk to LLMs (ChatGPT or Bard)
	#

	# select the model with -m
	# full model name supported, also several shorthands
	hof chat -m "gpt3" "why is the sky blue?" (gpt-3.5-turbo)
	hof chat -m "bard" "why is the sky blue?"  (chat-bison@001)

	# Ask of the LLM from strings, files, and/or stdin
	# these will be concatenated to from the question
	hof chat "Ask ChatGPT any question"    # as a string
	hof chat question.txt                  # from a file
	cat question.txt | hof chat -          # from stdin
	hof chat context.txt "and a question"  # mix all three

	# Provide a system message, these are special to LLMs
	# this is typically where the prompt engineering happens
	hof chat -S prompt.txt "now answer me this..."
	hof chat -S "... if short prompt ..." "now answer me this..."

	# Provide examples to the LLM
	# for Bard, these are an additional input
	# for ChatGPT, these will be appended to the system message
	# examples are supplied as JSON, they should be [{ input: string, output: string }]
	hof chat -E "<INPUT>: this is an input <OUTPUT>: this is an output" -E "..." "now answer me this..."
	hof chat -E examples.json "now answer me this"

	# Provide message history to the LLM
	# if messages are supplied as JSON, they should be { role: string, content: string }
	hof chat -M "user> asked some question" -M "assistant> had a reply" "now answer me this..."
	hof chat -M messages.json "now answer me this"


	"""#
