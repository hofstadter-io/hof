package functions

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var updateLong = `Update a Studios function by name with extra update values as input`

var UpdateCmd = &cobra.Command{

	Use: "update <name> <input>",

	Short: "Update a Studios function",

	Long: updateLong,

	Run: func(cmd *cobra.Command, args []string) {

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

		// Default body

		fmt.Println("hof studios functions update", name, input)

	},
}
