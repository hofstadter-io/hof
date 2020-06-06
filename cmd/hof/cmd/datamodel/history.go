package cmddatamodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/datamodel"
)

var historyLong = `show the history for a data model`

func HistoryRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = datamodel.RunHistoryFromArgs(args)

	return err
}

var HistoryCmd = &cobra.Command{

	Use: "history",

	Aliases: []string{
		"hist",
		"h",
		"log",
		"l",
	},

	Short: "show the history for a data model",

	Long: historyLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = HistoryRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := HistoryCmd.HelpFunc()
	usage := HistoryCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	HistoryCmd.SetHelpFunc(thelp)
	HistoryCmd.SetUsageFunc(tusage)

}