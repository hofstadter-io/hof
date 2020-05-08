package cmdconfig

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var testLong = `test your auth configuration, defaults to current`

func TestRun(name string) (err error) {

	return err
}

var TestCmd = &cobra.Command{

	Use: "test [name]",

	Short: "test your auth configuration, defaults to current",

	Long: testLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, strings.Join(args, "/"), 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

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
