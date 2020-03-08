package container

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/container"
)

var StatusLong = `Get the status of your container.
If name is not specified, the current directory is used.
`

var StatusCmd = &cobra.Command{

	Use: "status [name]",

	Short: "Get the status of your container",

	Long: StatusLong,

	Run: func(cmd *cobra.Command, args []string) {

		var name string
		if 0 < len(args) {
			name = args[0]
		}

		/*
			fmt.Println("hof containers status:",
				name,
			)
		*/

		err := container.Status(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
