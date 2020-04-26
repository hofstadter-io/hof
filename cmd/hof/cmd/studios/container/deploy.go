package cmdcontainer

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var deployLong = `Deploy a Studios container by name with extra update values as input`

func DeployRun(name string, input string) (err error) {

	return err
}

var DeployCmd = &cobra.Command{

	Use: "deploy <name> <input>",

	Short: "Deploy a Studios container",

	Long: deployLong,

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

		err = DeployRun(name, input)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
