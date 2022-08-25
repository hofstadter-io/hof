package cmdfmt

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var stopLong = `stop a formatter`

func StopRun(formatter string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var StopCmd = &cobra.Command{

	Use: "stop",

	Short: "stop a formatter",

	Long: stopLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'formatter'")
			cmd.Usage()
			os.Exit(1)
		}

		var formatter string

		if 0 < len(args) {

			formatter = args[0]

		}

		err = StopRun(formatter)
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

	ohelp := StopCmd.HelpFunc()
	ousage := StopCmd.UsageFunc()
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
	StopCmd.SetHelpFunc(thelp)
	StopCmd.SetUsageFunc(tusage)

}
