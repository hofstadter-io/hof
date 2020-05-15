package cmdmodel

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var testLong = `test the current migration and diff for a modelset`

func TestRun(name string) (err error) {

	return err
}

var TestCmd = &cobra.Command{

	Use: "test",

	Short: "test the current migration and diff for a modelset",

	Long: testLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'Name'")
			cmd.Usage()
			os.Exit(1)
		}

		var name string

		if 0 < len(args) {

			name = args[0]

		}

		err = TestRun(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
