package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

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
hof chat dm my-dm.cue "Add a Post model and make it so Users have many."`

func init() {

	ChatCmd.Flags().StringVarP(&(flags.ChatFlags.Model), "model", "M", "gpt-3.5-turbo", "LLM model to use [gpt-3.5-turbo,gpt-4]")
	ChatCmd.Flags().StringVarP(&(flags.ChatFlags.Prompt), "prompt", "P", "", "path to the system prompt, the first message in the chat")
	ChatCmd.Flags().StringVarP(&(flags.ChatFlags.Outfile), "outfile", "O", "", "path to write the output to")
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

}
