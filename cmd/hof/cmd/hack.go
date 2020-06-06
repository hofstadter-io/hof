package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/hack"
)

var hackLong = `development command`

func HackRun(args []string) (err error) {

	err = hack.Hack(args)

	// you can safely comment this print out
	// fmt.Println("not implemented")

	return err
}

var HackCmd = &cobra.Command{

	Use: "hack ...",

	Hidden: true,

	Short: "development command",

	Long: hackLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = HackRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := HackCmd.HelpFunc()
	usage := HackCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	HackCmd.SetHelpFunc(thelp)
	HackCmd.SetUsageFunc(tusage)

}