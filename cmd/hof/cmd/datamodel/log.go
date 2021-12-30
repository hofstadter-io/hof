package cmddatamodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var logLong = `show the current diff for a data model`

func LogRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var LogCmd = &cobra.Command{

	Use: "log",

	Aliases: []string{
		"l",
	},

	Short: "show the current diff for a data model",

	Long: logLong,

	PreRun: func(cmd *cobra.Command, args []string) {

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

	LogCmd.SetHelpFunc(help)
	LogCmd.SetUsageFunc(usage)

}
