package cmddatamodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/datamodel"
)

var listLong = `print available data models`

func ListRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = datamodel.RunListFromArgs(args, flags.DatamodelPflags)

	return err
}

var ListCmd = &cobra.Command{

	Use: "list",

	Aliases: []string{
		"ls",
	},

	Short: "print available data models",

	Long: listLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ListRun(args)
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

	ohelp := ListCmd.HelpFunc()
	ousage := ListCmd.UsageFunc()
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

	ListCmd.SetHelpFunc(help)
	ListCmd.SetUsageFunc(usage)

}
