package secret

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
)

var createLong = `Create a Studios secret from a file or key/val pairs`

var CreateCmd = &cobra.Command{

	Use: "create <name> <input>",

	Short: "Create a Studios secret",

	Long: createLong,

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

		fmt.Println("hof studios secret create", name, input)

	},
}
