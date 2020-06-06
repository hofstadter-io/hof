package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/resources"
)

var editLong = `edit resources`

func EditRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = resources.RunEditFromArgs(args)

	return err
}

var EditCmd = &cobra.Command{

	Use: "edit",

	Aliases: []string{
		"e",
	},

	Short: "edit resources",

	Long: editLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = EditRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := EditCmd.HelpFunc()
	usage := EditCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	EditCmd.SetHelpFunc(thelp)
	EditCmd.SetUsageFunc(tusage)

}