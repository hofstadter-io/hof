package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var rebaseLong = `Reapply commits on top of another base tip`

func RebaseRun(args []string) (err error) {

	return err
}

var RebaseCmd = &cobra.Command{

	Use: "rebase",

	Short: "Reapply commits on top of another base tip",

	Long: rebaseLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = RebaseRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
