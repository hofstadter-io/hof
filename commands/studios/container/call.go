package container

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/container"
)

var CallLong = `Call the container <name> with json <data>
data may be a JSON string or @filename.json
`

var CallCmd = &cobra.Command{

	Use: "call <name> [data]",

	Short: "Call a container with data",

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

		var data string
		if 1 < len(args) {
			data = args[1]
		}

		/*
			fmt.Println("hof containers call:",
				name,
				data,
			)
		*/

		err := container.Call(name, data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
