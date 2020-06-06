package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/learn"
)

var tutorialLong = `tutorials to help you learn hof right in hof`

func TutorialRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = learn.RunTutorialFromArgs(args)

	return err
}

var TutorialCmd = &cobra.Command{

	Use: "tutorial",

	Short: "tutorials to help you learn hof right in hof",

	Long: tutorialLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = TutorialRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := TutorialCmd.HelpFunc()
	usage := TutorialCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	TutorialCmd.SetHelpFunc(thelp)
	TutorialCmd.SetUsageFunc(tusage)

}