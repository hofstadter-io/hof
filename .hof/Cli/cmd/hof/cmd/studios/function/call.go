package cmdfunction

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var callLong = `Call your Studios function`

func CallRun(args []string) (err error) {

	return err
}

var CallCmd = &cobra.Command{

	Use: "call",

	Short: "Call a function",

	Long: callLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = CallRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
