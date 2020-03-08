package function

import (
	"fmt"
	"os"

	// custom imports

	// infered imports

	"github.com/hofstadter-io/hof/lib/fns"
	"github.com/spf13/cobra"
)

// Tool:   hof
// Name:   delete
// Usage:  delete <name or id>
// Parent: function

var DeleteLong = `Deletes the function <id> and all associated data in Hofstadter Studios.
You can get the id with "hof function list"
`

var DeleteCmd = &cobra.Command{

	Use: "delete <name or id>",

	Short: "Deletes the function by id",

	Long: DeleteLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In deleteCmd", "args", args)
		// Argument Parsing
		// [0]name:   name
		//     help:
		//     req'd:  false

		var name string

		if 0 < len(args) {

			name = args[0]
		}

		/*
			fmt.Println("hof function delete:",
				name,
			)
		*/

		err := fns.Delete(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {
	// add sub-commands to this command when present

}
