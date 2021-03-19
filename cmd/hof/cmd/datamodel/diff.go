package cmddatamodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/datamodel"
)

var diffLong = `show the current diff for a data model`

func DiffRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = datamodel.RunDiffFromArgs(args)

	return err
}

var DiffCmd = &cobra.Command{

	Use: "diff",

	Aliases: []string{
		"d",
	},

	Short: "show the current diff for a data model",

	Long: diffLong,

	PreRun: func(cmd *cobra.Command, args []string) {

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = DiffRun(args)
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

	ohelp := DiffCmd.HelpFunc()
	ousage := DiffCmd.UsageFunc()
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

	DiffCmd.SetHelpFunc(help)
	DiffCmd.SetUsageFunc(usage)

}
