package cmdmodel

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var createLong = `create a modelset`

func CreateRun(name string, entrypoint string) (err error) {

	return err
}

var CreateCmd = &cobra.Command{

	Use: "create",

	Short: "create a modelset",

	Long: createLong,

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

		var entrypoint string
		entrypoint = "models"

		if 1 < len(args) {

			entrypoint = args[1]

		}

		err = CreateRun(name, entrypoint)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
