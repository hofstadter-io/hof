package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib"
)

var feedbackLong = `send feedback, bug reports, or any message :]
	email:     (optional) your email, if you'd like us to reply
	message:   your message, please be respectful to the person receiving it`

func FeedbackRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = lib.SendFeedback(args)

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

	Short: "send feedback, bug reports, or any message :]",

	Long: feedbackLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "", 0)

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

	help := FeedbackCmd.HelpFunc()
	usage := FeedbackCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "", 0)
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/usage", "", 0)
		return usage(cmd)
	}
	FeedbackCmd.SetHelpFunc(thelp)
	FeedbackCmd.SetUsageFunc(tusage)

}