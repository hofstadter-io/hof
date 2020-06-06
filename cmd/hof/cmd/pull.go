package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/workspace"
)

var pullLong = `fetch from and integrate with another repository or a local branch`

func PullRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = workspace.RunPullFromArgs(args)

	return err
}

var PullCmd = &cobra.Command{

	Use: "pull",

	Short: "fetch from and integrate with another repository or a local branch",

	Long: pullLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
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

	help := PullCmd.HelpFunc()
	usage := PullCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	PullCmd.SetHelpFunc(thelp)
	PullCmd.SetUsageFunc(tusage)

}