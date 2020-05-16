package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var checkoutLong = `Switch branches or restore working tree files`

func CheckoutRun(args []string) (err error) {

	return err
}

var CheckoutCmd = &cobra.Command{

	Use: "checkout",

	Short: "Switch branches or restore working tree files",

	Long: checkoutLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

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
