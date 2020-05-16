package cmd

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var pushLong = ` Update remote refs along with associated objects`

func PushRun(args []string) (err error) {

	return err
}

var PushCmd = &cobra.Command{

	Use: "push",

	Short: " Update remote refs along with associated objects",

	Long: pushLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = PushRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
