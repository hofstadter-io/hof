package cmddatabase

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var updateLong = `Update a Studios database by name with extra update values as input`

func UpdateRun(name string, input string) (err error) {

	return err
}

var UpdateCmd = &cobra.Command{

	Use: "update <name> <input>",

	Short: "Update a Studios database",

	Long: updateLong,

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
			fmt.Println("missing required argument: 'Input'")
			cmd.Usage()
			os.Exit(1)
		}

		var input string

		if 1 < len(args) {

			input = args[1]

		}

		err = UpdateRun(name, input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
