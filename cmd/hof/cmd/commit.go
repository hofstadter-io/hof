package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/workspace"
)

var commitLong = `record changes to the repository`

func CommitRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = workspace.RunCommitFromArgs(args)

	return err
}

var CommitCmd = &cobra.Command{

	Use: "commit",

	Short: "record changes to the repository",

	Long: commitLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = CommitRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := CommitCmd.HelpFunc()
	usage := CommitCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	CommitCmd.SetHelpFunc(thelp)
	CommitCmd.SetUsageFunc(tusage)

}