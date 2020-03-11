package container

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/container"
)

var LogsLong = `List the logs of your container.
If name is not specified, the current directory is used.
`

var LogsCmd = &cobra.Command{

	Use: "logs [name]",

	Short: "List the logs of your container",

	Long: LogsLong,

	Run: func(cmd *cobra.Command, args []string) {

		var name string
		if 0 < len(args) {
			name = args[0]
		}

		/*
			fmt.Println("hof containers logs:",
				name,
			)
		*/

		err := container.Logs(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
