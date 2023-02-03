package cmddatamodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/datamodel/cmd"
)

var logLong = `show the history for a datamodel`

func init() {

	LogCmd.Flags().BoolVarP(&(flags.Datamodel__LogFlags.ByValue), "by-value", "", false, "display snapshot log by value")
	LogCmd.Flags().BoolVarP(&(flags.Datamodel__LogFlags.Details), "details", "", false, "print more when displaying the log")
}

func LogRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = cmd.Run("log", args, flags.RootPflags, flags.DatamodelPflags)

	return err
}

var LogCmd = &cobra.Command{

	Use: "log",

	Aliases: []string{
		"l",
	},

	Short: "show the history for a datamodel",

	Long: logLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = LogRun(args)
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

	ohelp := LogCmd.HelpFunc()
	ousage := LogCmd.UsageFunc()
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
	LogCmd.SetHelpFunc(thelp)
	LogCmd.SetUsageFunc(tusage)

}
