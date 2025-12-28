package cmdenv

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var pullLong = `pull an environment`

func PullRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var PullCmd = &cobra.Command{

	Use: "pull [...target] [% ...cue]",

	Short: "pull an environment",

	Long: pullLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = PullRun(args)
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

	ohelp := PullCmd.HelpFunc()
	ousage := PullCmd.UsageFunc()

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
	PullCmd.SetHelpFunc(thelp)
	PullCmd.SetUsageFunc(tusage)

}
