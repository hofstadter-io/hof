package cmdconfig

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var setLong = `set config values with an expr`

func SetRun(expr string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var SetCmd = &cobra.Command{

	Use: "set [expr]",

	Short: "set config values with an expr",

	Long: setLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'expr'")
			cmd.Usage()
			os.Exit(1)
		}

		var expr string

		if 0 < len(args) {

			expr = args[0]

		}

		err = SetRun(expr)
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

	ohelp := SetCmd.HelpFunc()
	ousage := SetCmd.UsageFunc()
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
	SetCmd.SetHelpFunc(thelp)
	SetCmd.SetUsageFunc(tusage)

}
