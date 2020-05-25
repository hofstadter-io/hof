package cmddatamodel

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/datamodel"
)

var diffLong = `show the current diff for a data model`

func DiffRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = datamodel.RunDiffFromArgs(args)

	return err
}

var DiffCmd = &cobra.Command{

	Use: "diff",

	Aliases: []string{
		"d",
	},

	Short: "show the current diff for a data model",

	Long: diffLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = DiffRun(args)
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
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/usage", "<omit>", 0)
		return usage(cmd)
	}
	DiffCmd.SetHelpFunc(thelp)
	DiffCmd.SetUsageFunc(tusage)

}