package function

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/function"
)

var DeleteLong = `Deletes the function <id> and all associated data in Hofstadter Studios.
You can get the id with "hof function list"
`

var DeleteCmd = &cobra.Command{

	Use: "delete <name or id>",

	Short: "Deletes the function by id",

	Long: DeleteLong,

	Run: func(cmd *cobra.Command, args []string) {

		var name string
		if 0 < len(args) {
			name = args[0]
		}

		/*
			fmt.Println("hof function delete:",
				name,
			)
		*/

		err := function.Delete(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
