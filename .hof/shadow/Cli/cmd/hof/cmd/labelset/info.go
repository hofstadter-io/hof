package cmdlabelset

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var infoLong = `print info about labelsets in your workspace or system`

func InfoRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var InfoCmd = &cobra.Command{

	Use: "info",

	Aliases: []string{
		"i",
	},

	Short: "print info about labelsets in your workspace or system",

	Long: infoLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = InfoRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := InfoCmd.HelpFunc()
	usage := InfoCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	InfoCmd.SetHelpFunc(thelp)
	InfoCmd.SetUsageFunc(tusage)

}
