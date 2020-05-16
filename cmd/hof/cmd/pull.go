package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var pullLong = `Fetch from and integrate with another repository or a local branch`

func PullRun(args []string) (err error) {

	return err
}

var PullCmd = &cobra.Command{

	Use: "pull",

	Short: "Fetch from and integrate with another repository or a local branch",

	Long: pullLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = PullRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
