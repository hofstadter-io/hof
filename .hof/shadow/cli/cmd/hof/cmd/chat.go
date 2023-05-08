package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var chatLong = `hof chat synergizes LLM technology with our code gen technology.
You can chat with the AI to help you write designs and code.
Module authors can provide custom prompts for their schemas.

Currently, only ChatGPT is supported. You can use any of the
gpt-3.5 or gpt-4 models. The flag should match OpenAI API options.

Set OPENAI_API_KEY`

func init() {

	ChatCmd.Flags().StringVarP(&(flags.ChatFlags.Model), "model", "M", "gpt-3.5-turbo", "LLM model to use [gpt-3.5-turbo,gpt-4]")
	ChatCmd.Flags().StringVarP(&(flags.ChatFlags.Parameters), "parameters", "P", "", "path to a config file containing model parameters")
	ChatCmd.Flags().StringSliceVarP(&(flags.ChatFlags.Generator), "generator", "G", nil, "generator tags to run, default is all")
}

func ChatRun(jsonfile string, instructions string, extra []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var ChatCmd = &cobra.Command{

	Use: "chat jsonfile [extra args]",

	Short: "Co-design with AI (alpha)",

	Long: chatLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'jsonfile'")
			cmd.Usage()
			os.Exit(1)
		}

		var jsonfile string

		if 0 < len(args) {

			jsonfile = args[0]

		}

		if 1 >= len(args) {
			fmt.Println("missing required argument: 'instructions'")
			cmd.Usage()
			os.Exit(1)
		}

		var instructions string

		if 1 < len(args) {

			instructions = args[1]

		}

		var extra []string

		if 2 < len(args) {

			extra = args[2:]

		}

		err = ChatRun(jsonfile, instructions, extra)
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
