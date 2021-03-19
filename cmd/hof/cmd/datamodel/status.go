package cmddatamodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/datamodel"
)

var statusLong = `print the data model status`

func StatusRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = datamodel.RunStatusFromArgs(args)

	return err
}

var StatusCmd = &cobra.Command{

	Use: "status",

	Aliases: []string{
		"st",
	},

	Short: "print the data model status",

	Long: statusLong,

	PreRun: func(cmd *cobra.Command, args []string) {

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = StatusRun(args)
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

	ohelp := StatusCmd.HelpFunc()
	ousage := StatusCmd.UsageFunc()
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

	StatusCmd.SetHelpFunc(help)
	StatusCmd.SetUsageFunc(usage)

}
