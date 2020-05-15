package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var uiLong = `run hof's local web ui`

func UiRun(args []string) (err error) {

	return err
}

var UiCmd = &cobra.Command{

	Use: "ui",

	Short: "run hof's local web ui",

	Long: uiLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = UiRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
