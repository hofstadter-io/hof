package cmdapp

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var statusLong = `Get the status of a Studios app.`

func StatusRun(ident string) (err error) {

	return err
}

var StatusCmd = &cobra.Command{

	Use: "status <name or id>",

	Short: "Get the status of a Studios app.",

	Long: statusLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, strings.Join(args, "/"), 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'Ident'")
			cmd.Usage()
			os.Exit(1)
		}

		var ident string

		if 0 < len(args) {

			ident = args[0]

		}

		err = StatusRun(ident)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
