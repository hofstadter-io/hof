package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib"
)

var feedbackLong = `Opens an issue or discussion on GitHub with some fields prefilled out`

func init() {

	flags.SetupFeedbackPflags(FeedbackCmd.PersistentFlags(), &(flags.FeedbackPflags))

}

func FeedbackRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = lib.SendFeedback(args, flags.RootPflags, flags.FeedbackPflags)

	return err
}

var FeedbackCmd = &cobra.Command{

	Use: "feedback <message>",

	Aliases: []string{
		"hi",
		"ask",
		"report",
	},

	Short: "open an issue or discussion on GitHub",

	Long: feedbackLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

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
	FeedbackCmd.SetHelpFunc(thelp)
	FeedbackCmd.SetUsageFunc(tusage)

}
