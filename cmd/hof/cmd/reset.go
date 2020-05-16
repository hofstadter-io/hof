package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var resetLong = `Reset current HEAD to the specified state`

func ResetRun(args []string) (err error) {

	return err
}

var ResetCmd = &cobra.Command{

	Use: "reset",

	Short: "Reset current HEAD to the specified state",

	Long: resetLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ResetRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
