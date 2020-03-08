package function

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/function"
)

var CallLong = `Call the function <name> with <data>
data may be a JSON string or @filename.json
`

var CallCmd = &cobra.Command{

	Use: "call <name> <data>",

	Short: "Call a function by name",

	Long: CallLong,

	Run: func(cmd *cobra.Command, args []string) {

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'name'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var name string
		if 0 < len(args) {
			name = args[0]
		}

		if 1 >= len(args) {
			fmt.Println("missing required argument: 'data'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var data string
		if 1 < len(args) {
			data = args[1]
		}

		/*
			fmt.Println("hof function call:",
				name,
				data,
			)
		*/

		err := function.Call(name, data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
