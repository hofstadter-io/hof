package cmd

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/chat"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var chatLong = `Use chat to work with hof features or from modules you import.
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

`

func init() {

	flags.SetupChatFlags(ChatCmd.Flags(), &(flags.ChatFlags))

}

func ChatRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var ChatCmd = &cobra.Command{

	Use: "chat [args]",

	Short: "co-create with AI (alpha)",

	Long: chatLong,

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		glob := toComplete + "*"
		matches, _ := filepath.Glob(glob)
		return matches, cobra.ShellCompDirectiveDefault
	},

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = ChatRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := ChatCmd.HelpFunc()
	ousage := ChatCmd.UsageFunc()

	help := func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath() + " help")

		if extra(cmd) {
			return
		}
		ohelp(cmd, args)
	}
	usage := func(cmd *cobra.Command) error {
		if extra(cmd) {
			return nil
		}
		return ousage(cmd)
	}

	thelp := func(cmd *cobra.Command, args []string) {
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		return usage(cmd)
	}
	ChatCmd.SetHelpFunc(thelp)
	ChatCmd.SetUsageFunc(tusage)

	ChatCmd.AddCommand(cmdchat.ListCmd)
	ChatCmd.AddCommand(cmdchat.WithCmd)
	ChatCmd.AddCommand(cmdchat.InfoCmd)

}
