package cmdstore

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var runLong = `run local datastore servers`

func RunRun(dstype string, name string) (err error) {

	return err
}

var RunCmd = &cobra.Command{

	Use: "run",

	Aliases: []string{
		"local",
	},

	Short: "run local datastore servers",

	Long: runLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'Dstype'")
			cmd.Usage()
			os.Exit(1)
		}

		var dstype string

		if 0 < len(args) {

			dstype = args[0]

		}

		if 1 >= len(args) {
			fmt.Println("missing required argument: 'Name'")
			cmd.Usage()
			os.Exit(1)
		}

		var name string

		if 1 < len(args) {

			name = args[1]

		}

		err = RunRun(dstype, name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
