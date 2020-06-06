package cmdst

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/structural"
)

var diffLong = `Calculate the difference between two Cue values`

func DiffRun(orig string, next string, entrypoints []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = structural.RunDiffFromArgs(orig, next, entrypoints)

	return err
}

var DiffCmd = &cobra.Command{

	Use: "diff <orig> <next> [...entrypoints]",

	Short: "calculate the difference between two Cue values",

	Long: diffLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'orig'")
			cmd.Usage()
			os.Exit(1)
		}

		var orig string

		if 0 < len(args) {

			orig = args[0]

		}

		if 1 >= len(args) {
			fmt.Println("missing required argument: 'next'")
			cmd.Usage()
			os.Exit(1)
		}

		var next string

		if 1 < len(args) {

			next = args[1]

		}

		var entrypoints []string

		if 2 < len(args) {

			entrypoints = args[2:]

		}

		err = DiffRun(orig, next, entrypoints)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := DiffCmd.HelpFunc()
	usage := DiffCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	DiffCmd.SetHelpFunc(thelp)
	DiffCmd.SetUsageFunc(tusage)

}