package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/workspace"
)

var checkoutLong = `switch branches or restore working tree files`

func CheckoutRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = workspace.RunCheckoutFromArgs(args)

	return err
}

var CheckoutCmd = &cobra.Command{

	Use: "checkout",

	Short: "switch branches or restore working tree files",

	Long: checkoutLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = CheckoutRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := CheckoutCmd.HelpFunc()
	usage := CheckoutCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	CheckoutCmd.SetHelpFunc(thelp)
	CheckoutCmd.SetUsageFunc(tusage)

}