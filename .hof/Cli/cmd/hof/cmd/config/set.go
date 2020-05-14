package cmdconfig

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var setLong = `set configuration values`

func SetRun(name string, host string, account string, project string) (err error) {

	return err
}

var SetCmd = &cobra.Command{

	Use: "set <name> <host> <account> [project]",

	Short: "set configuration values",

	Long: setLong,

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

		if 1 >= len(args) {
			fmt.Println("missing required argument: 'Host'")
			cmd.Usage()
			os.Exit(1)
		}

		var host string

		if 1 < len(args) {

			host = args[1]

		}

		if 2 >= len(args) {
			fmt.Println("missing required argument: 'Account'")
			cmd.Usage()
			os.Exit(1)
		}

		var account string

		if 2 < len(args) {

			account = args[2]

		}

		var project string

		if 3 < len(args) {

			project = args[3]

		}

		err = SetRun(name, host, account, project)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
