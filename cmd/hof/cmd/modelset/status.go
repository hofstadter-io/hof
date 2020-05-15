package cmdmodelset

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var statusLong = `show the current status for a modelset`

func StatusRun(name string) (err error) {

	return err
}

var StatusCmd = &cobra.Command{

	Use: "status",

	Short: "show the current status for a modelset",

	Long: statusLong,

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

		err = StatusRun(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
