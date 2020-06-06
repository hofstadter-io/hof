package cmdlabel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/labels"
)

var applyLong = `find and apply labels to resources`

func ApplyRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = labels.RunApplyLabelFromArgs(args)

	return err
}

var ApplyCmd = &cobra.Command{

	Use: "apply",

	Aliases: []string{
		"a",
	},

	Short: "find and apply labels to resources",

	Long: applyLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ApplyRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := ApplyCmd.HelpFunc()
	usage := ApplyCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	ApplyCmd.SetHelpFunc(thelp)
	ApplyCmd.SetUsageFunc(tusage)

}