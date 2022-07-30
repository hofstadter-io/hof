package cmdfmt

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var startLong = `start a formatter`

func StartRun(formatter string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var StartCmd = &cobra.Command{

	Use: "start",

	Short: "start a formatter",

	Long: startLong,

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

		err = StartRun(formatter)
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

	ohelp := StartCmd.HelpFunc()
	ousage := StartCmd.UsageFunc()
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

	StartCmd.SetHelpFunc(help)
	StartCmd.SetUsageFunc(usage)

}
