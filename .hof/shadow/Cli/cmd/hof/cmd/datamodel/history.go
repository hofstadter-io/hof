package cmddatamodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var historyLong = `show the history for a data model`

func HistoryRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

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
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := HistoryCmd.HelpFunc()
	ousage := HistoryCmd.UsageFunc()
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

	HistoryCmd.SetHelpFunc(help)
	HistoryCmd.SetUsageFunc(usage)

}
