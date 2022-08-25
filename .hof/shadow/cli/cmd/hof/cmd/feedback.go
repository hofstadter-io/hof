package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var feedbackLong = `send feedback, bug reports, or any message
	email:     (optional) your email, if you'd like us to reply
	message:   your message, please be respectful to the person receiving it`

func FeedbackRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var FeedbackCmd = &cobra.Command{

	Use: "feedback [email] <message>",

	Aliases: []string{
		"hi",
		"say",
		"from",
		"bug",
		"yo",
		"hello",
		"greetings",
		"support",
	},

	Short: "send feedback, bug reports, or any message",

	Long: feedbackLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = FeedbackRun(args)
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

	ohelp := FeedbackCmd.HelpFunc()
	ousage := FeedbackCmd.UsageFunc()
	help := func(cmd *cobra.Command, args []string) {
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
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	FeedbackCmd.SetHelpFunc(thelp)
	FeedbackCmd.SetUsageFunc(tusage)

}
