package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/ops"
)

var runLong = `run polyglot command and scripts seamlessly across runtimes`

func RunRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = ops.RunRunFromArgs(args)

	return err
}

var RunCmd = &cobra.Command{

	Use: "run",

	Aliases: []string{
		"r",
	},

	Short: "run polyglot command and scripts seamlessly across runtimes",

	Long: runLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = RunRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := RunCmd.HelpFunc()
	usage := RunCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	RunCmd.SetHelpFunc(thelp)
	RunCmd.SetUsageFunc(tusage)

}