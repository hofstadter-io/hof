package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var devLong = `run hof's local dev setup`

func DevRun(args []string) (err error) {

	return err
}

var DevCmd = &cobra.Command{

	Use: "dev",

	Short: "run hof's local dev setup",

	Long: devLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = DevRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
