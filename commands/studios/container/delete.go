package container

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/container"
)

var DeleteLong = `Deletes a container by <id> and all associated data in Studios.
You can find the id by running 'hof cont list'
`

var DeleteCmd = &cobra.Command{

	Use: "delete <name or id>",

	Short: "Deletes a container",

	Long: DeleteLong,

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

		/*
			fmt.Println("hof containers delete:",
				id,
			)
		*/

		err := container.Delete(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
