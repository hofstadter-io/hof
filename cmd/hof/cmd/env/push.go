package cmdenv

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/cmd/hof/ga"
	"github.com/hofstadter-io/hof/lib/env/cmd"
)

var pushLong = `push an environment`

func PushRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = cmd.Push(args, flags.RootPflags)

	return err
}

var PushCmd = &cobra.Command{

	Use: "push <name>",

	Short: "push an environment",

	Long: pushLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = PushRun(args)
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

	ohelp := PushCmd.HelpFunc()
	ousage := PushCmd.UsageFunc()

	help := func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath() + " help")

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
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		return usage(cmd)
	}
	PushCmd.SetHelpFunc(thelp)
	PushCmd.SetUsageFunc(tusage)

}
