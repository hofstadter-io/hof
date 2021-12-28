package cmddatamodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listLong = `find and display data models`

func ListRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var ListCmd = &cobra.Command{

	Use: "list",

	Aliases: []string{
		"l",
	},

	Short: "find and display data models",

	Long: listLong,

	PreRun: func(cmd *cobra.Command, args []string) {

	},

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
