package container

import (
	"fmt"

	// custom imports

	// infered imports
	"os"

	"github.com/hofstadter-io/hof/lib/crun"
	"github.com/spf13/cobra"
)

// Tool:   hof
// Name:   delete
// Usage:  delete <name or id>
// Parent: container

var DeleteLong = `Deletes a container by <id> and all associated data in Studios.
You can find the id by running 'hof crun list'
`

var DeleteCmd = &cobra.Command{

	Use: "delete <name or id>",

	Short: "Deletes a container",

	Long: DeleteLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In deleteCmd", "args", args)
		// Argument Parsing
		// [0]name:   name
		//     help:
		//     req'd:  true
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

		err := crun.Delete(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	// add sub-commands to this command when present

}
