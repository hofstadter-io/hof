package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/ops"
)

var addLong = `add dependencies and new components to the current module or workspace`

func AddRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = ops.RunAddFromArgs(args)

	return err
}

var AddCmd = &cobra.Command{

	Use: "add",

	Short: "add dependencies and new components to the current module or workspace",

	Long: addLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = AddRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := AddCmd.HelpFunc()
	usage := AddCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	AddCmd.SetHelpFunc(thelp)
	AddCmd.SetUsageFunc(tusage)

}