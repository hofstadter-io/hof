package cmdmod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/mod"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var graphLong = `print module requirement graph`

func GraphRun(args []string) (err error) {

	err = mod.ProcessLangs("graph", args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return err
}

var GraphCmd = &cobra.Command{

	Use: "graph",

	Short: "print module requirement graph",

	Long: graphLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = GraphRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := GraphCmd.HelpFunc()
	usage := GraphCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	GraphCmd.SetHelpFunc(thelp)
	GraphCmd.SetUsageFunc(tusage)

}
