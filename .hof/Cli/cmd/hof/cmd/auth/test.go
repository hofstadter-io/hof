package cmdauth

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var testLong = `test your auth configuration, defaults to current`

func TestRun(args []string) (err error) {

	return err
}

var TestCmd = &cobra.Command{

	Use: "test [name]",

	Short: "test your auth configuration, defaults to current",

	Long: testLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = TestRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
