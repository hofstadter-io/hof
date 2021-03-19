package cmddatamodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var applyLong = `apply a migraion sequence against a data store`

func ApplyRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var ApplyCmd = &cobra.Command{

	Use: "apply",

	Aliases: []string{
		"a",
	},

	Short: "apply a migraion sequence against a data store",

	Long: applyLong,

	PreRun: func(cmd *cobra.Command, args []string) {

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ApplyRun(args)
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

	ohelp := ApplyCmd.HelpFunc()
	ousage := ApplyCmd.UsageFunc()
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

	ApplyCmd.SetHelpFunc(help)
	ApplyCmd.SetUsageFunc(usage)

}
